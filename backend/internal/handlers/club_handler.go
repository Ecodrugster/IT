package handlers

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

type ClubInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

func normalizeClubInput(input *ClubInput) bool {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.Icon = strings.TrimSpace(input.Icon)
	input.Color = strings.TrimSpace(input.Color)

	if input.Name == "" || input.Description == "" {
		return false
	}

	if input.Icon == "" {
		input.Icon = "👥"
	}
	if input.Color == "" {
		input.Color = "bg-blue-600/20"
	}

	return true
}

func getClubStatus(data map[string]interface{}) string {
	status := asString(data["status"])
	if status == "" {
		return "approved"
	}
	return status
}

func GetClubs(c *gin.Context) {
	uid, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, err := getUserRoleByUID(c, uid)
	if err != nil {
		role = "student"
	}
	canSeeAll := isAdminLikeRole(role)

	iter := repositories.FirestoreClient.Collection("clubs").Documents(c.Request.Context())

	var clubs []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clubs"})
			return
		}

		data := doc.Data()
		status := getClubStatus(data)
		createdBy := asString(data["created_by"])

		if !canSeeAll && status != "approved" && createdBy != uid {
			continue
		}

		data["id"] = doc.Ref.ID
		data["status"] = status
		clubs = append(clubs, data)
	}

	sort.Slice(clubs, func(i, j int) bool {
		return asString(clubs[i]["name"]) < asString(clubs[j]["name"])
	})

	c.JSON(http.StatusOK, clubs)
}

func CreateClub(c *gin.Context) {
	uid, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ClubInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !normalizeClubInput(&input) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and description are required"})
		return
	}

	now := time.Now()
	club := map[string]interface{}{
		"name":         input.Name,
		"description":  input.Description,
		"icon":         input.Icon,
		"color":        input.Color,
		"members":      []string{uid},
		"created_by":   uid,
		"status":       "pending",
		"created_at":   now,
		"updated_at":   now,
		"moderated_at": nil,
		"moderated_by": "",
	}

	_, _, err = repositories.FirestoreClient.Collection("clubs").Add(c.Request.Context(), club)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create club request"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Club request sent for moderation",
		"status":  "pending",
	})
}

func JoinClub(c *gin.Context) {
	clubID := c.Param("id")
	userID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ref := repositories.FirestoreClient.Collection("clubs").Doc(clubID)
	doc, err := ref.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	if getClubStatus(doc.Data()) != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Club is not approved yet"})
		return
	}

	_, err = ref.Update(c.Request.Context(), []firestore.Update{
		{
			Path:  "members",
			Value: firestore.ArrayUnion(userID),
		},
		{
			Path:  "updated_at",
			Value: time.Now(),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join club"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined successfully"})
}

func UpdateClub(c *gin.Context) {
	clubID := c.Param("id")
	userID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ref := repositories.FirestoreClient.Collection("clubs").Doc(clubID)
	doc, err := ref.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	if asString(doc.Data()["created_by"]) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the creator can edit this club"})
		return
	}

	var input ClubInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if !normalizeClubInput(&input) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and description are required"})
		return
	}

	current := doc.Data()
	currentStatus := getClubStatus(current)
	needsModeration :=
		asString(current["name"]) != input.Name ||
			asString(current["description"]) != input.Description ||
			asString(current["icon"]) != input.Icon ||
			asString(current["color"]) != input.Color

	update := map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
		"icon":        input.Icon,
		"color":       input.Color,
		"updated_at":  time.Now(),
	}

	if needsModeration && (currentStatus == "approved" || currentStatus == "rejected") {
		update["status"] = "pending"
		update["moderated_by"] = ""
		update["moderated_at"] = nil
		update["moderation_comment"] = ""
	}

	_, err = ref.Set(c.Request.Context(), update, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update club"})
		return
	}

	message := "Club updated successfully"
	if status, ok := update["status"].(string); ok && status == "pending" {
		message = "Club updated and sent for re-moderation"
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func DeleteClub(c *gin.Context) {
	clubID := c.Param("id")
	userID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ref := repositories.FirestoreClient.Collection("clubs").Doc(clubID)
	doc, err := ref.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	if asString(doc.Data()["created_by"]) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the creator can delete this club"})
		return
	}

	_, err = ref.Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete club"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Club deleted successfully"})
}

func LeaveClub(c *gin.Context) {
	clubID := c.Param("id")
	userID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ref := repositories.FirestoreClient.Collection("clubs").Doc(clubID)
	_, err = ref.Update(c.Request.Context(), []firestore.Update{
		{
			Path:  "members",
			Value: firestore.ArrayRemove(userID),
		},
		{
			Path:  "updated_at",
			Value: time.Now(),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave club"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left successfully"})
}
