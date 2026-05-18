package handlers

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

// AdminGetUsers returns users list with basic pagination.
func AdminGetUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit := 50
	fmt.Sscanf(limitStr, "%d", &limit)

	iter := repositories.FirestoreClient.Collection("users").Limit(limit).Documents(c.Request.Context())

	var users []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		users = append(users, normalizeUserDocument(doc.Data(), doc.Ref.ID))
	}

	c.JSON(http.StatusOK, gin.H{"users": users, "total": len(users)})
}

// AdminUpdateUserRole changes user role.
func AdminUpdateUserRole(c *gin.Context) {
	uid := strings.TrimSpace(c.Param("id"))
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User id is required"})
		return
	}
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	input.Role = strings.TrimSpace(input.Role)

	if input.Role != "admin" && input.Role != "student" && input.Role != "teacher" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role value"})
		return
	}

	oldRole := "student"
	targetName := uid
	if userDoc, getErr := repositories.FirestoreClient.Collection("users").Doc(uid).Get(c.Request.Context()); getErr == nil {
		oldRole = asString(userDoc.Data()["role"])
		if oldRole == "" {
			oldRole = "student"
		}
		targetName = getDisplayNameFromUserData(userDoc.Data(), uid)
	}

	_, err = repositories.FirestoreClient.Collection("users").Doc(uid).Update(c.Request.Context(), []firestore.Update{{Path: "role", Value: input.Role}})
	if err != nil {
		_, err = repositories.FirestoreClient.Collection("users").Doc(uid).Set(c.Request.Context(), map[string]interface{}{
			"role":       input.Role,
			"updated_at": time.Now(),
		}, repositories.MergeAll)
	} else {
		_, _ = repositories.FirestoreClient.Collection("users").Doc(uid).Set(c.Request.Context(), map[string]interface{}{
			"updated_at": time.Now(),
		}, repositories.MergeAll)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	logAdminAction(c, adminID, "user.role.updated", "user", uid, targetName, map[string]interface{}{
		"old_role": oldRole,
		"new_role": input.Role,
	})

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

// AdminUpdateUserGroup assigns or clears student's group.
func AdminUpdateUserGroup(c *gin.Context) {
	uid := strings.TrimSpace(c.Param("id"))
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User id is required"})
		return
	}

	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		GroupName string `json:"group_name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	groupName := strings.TrimSpace(input.GroupName)
	targetName := uid
	oldGroup := ""
	if userDoc, getErr := repositories.FirestoreClient.Collection("users").Doc(uid).Get(c.Request.Context()); getErr == nil {
		targetName = getDisplayNameFromUserData(userDoc.Data(), uid)
		oldGroup = asString(userDoc.Data()["group_name"])
		if oldGroup == "" {
			oldGroup = asString(userDoc.Data()["group"])
		}
	}

	update := map[string]interface{}{
		"group_name": groupName,
		"group":      groupName,
		"updated_at": time.Now(),
	}

	_, err = repositories.FirestoreClient.Collection("users").Doc(uid).Set(c.Request.Context(), update, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user group"})
		return
	}

	logAdminAction(c, adminID, "user.group.updated", "user", uid, targetName, map[string]interface{}{
		"old_group": oldGroup,
		"new_group": groupName,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":    "User group updated successfully",
		"group_name": groupName,
	})
}

// BootstrapFirstAdmin assigns admin role to the current user only if no admins exist yet.
func BootstrapFirstAdmin(c *gin.Context) {
	uid, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentRole, _ := getUserRoleByUID(c, uid)
	if currentRole == "admin" {
		c.JSON(http.StatusOK, gin.H{"message": "User is already admin"})
		return
	}

	adminIter := repositories.FirestoreClient.Collection("users").Where("role", "==", "admin").Limit(1).Documents(c.Request.Context())
	_, err = adminIter.Next()
	if err != iterator.Done && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing admins"})
		return
	}
	if err != iterator.Done {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin already exists. Ask current admin to change your role."})
		return
	}

	_, err = repositories.FirestoreClient.Collection("users").Doc(uid).Set(c.Request.Context(), map[string]interface{}{
		"role":       "admin",
		"updated_at": time.Now(),
	}, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign admin role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin role assigned successfully", "uid": uid})
}

// AdminGetPosts returns recent posts for moderation.
func AdminGetPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit := 50
	fmt.Sscanf(limitStr, "%d", &limit)

	iter := repositories.FirestoreClient.Collection("posts").OrderBy("created_at", repositories.Descending).Limit(limit).Documents(c.Request.Context())

	var posts []map[string]interface{}
	nameCache := map[string]string{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
			return
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		authorUID := asString(data["author_id"])
		if authorUID != "" {
			data["author_name"] = getUserDisplayNameByUID(c, authorUID, nameCache)
		}
		posts = append(posts, data)
	}

	c.JSON(http.StatusOK, posts)
}

// AdminDeletePost removes a post.
func AdminDeletePost(c *gin.Context) {
	id := c.Param("id")
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	targetName := id
	details := map[string]interface{}{}
	if doc, getErr := repositories.FirestoreClient.Collection("posts").Doc(id).Get(c.Request.Context()); getErr == nil {
		post := doc.Data()
		if content := asString(post["content"]); content != "" {
			if len(content) > 80 {
				content = content[:80] + "..."
			}
			targetName = content
		}
		if authorUID := asString(post["author_id"]); authorUID != "" {
			details["author_uid"] = authorUID
			details["author_name"] = getUserDisplayNameByUID(c, authorUID, nil)
		}
	}

	_, err = repositories.FirestoreClient.Collection("posts").Doc(id).Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	logAdminAction(c, adminID, "post.deleted", "post", id, targetName, details)

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted by moderator"})
}

// AdminDeleteNews removes news item.
func AdminDeleteNews(c *gin.Context) {
	id := c.Param("id")
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	targetName := id
	if doc, getErr := repositories.FirestoreClient.Collection("news").Doc(id).Get(c.Request.Context()); getErr == nil {
		title := asString(doc.Data()["title"])
		if title != "" {
			targetName = title
		}
	}

	_, err = repositories.FirestoreClient.Collection("news").Doc(id).Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news"})
		return
	}

	logAdminAction(c, adminID, "news.deleted", "news", id, targetName, nil)
	c.JSON(http.StatusOK, gin.H{"message": "News deleted"})
}

// AdminUpdateNews edits news item.
func AdminUpdateNews(c *gin.Context) {
	id := c.Param("id")
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	_, err = repositories.FirestoreClient.Collection("news").Doc(id).Set(c.Request.Context(), input, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news"})
		return
	}

	targetName := asString(input["title"])
	if targetName == "" {
		targetName = id
	}
	logAdminAction(c, adminID, "news.updated", "news", id, targetName, nil)

	c.JSON(http.StatusOK, gin.H{"message": "News updated"})
}

func readAllClubs(c *gin.Context) ([]map[string]interface{}, error) {
	iter := repositories.FirestoreClient.Collection("clubs").Documents(c.Request.Context())
	var clubs []map[string]interface{}
	nameCache := map[string]string{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		data["status"] = getClubStatus(data)
		if createdBy := asString(data["created_by"]); createdBy != "" {
			data["created_by_name"] = getUserDisplayNameByUID(c, createdBy, nameCache)
		}
		if moderatedBy := asString(data["moderated_by"]); moderatedBy != "" {
			data["moderated_by_name"] = getUserDisplayNameByUID(c, moderatedBy, nameCache)
		}
		clubs = append(clubs, data)
	}

	sort.Slice(clubs, func(i, j int) bool {
		statusI := asString(clubs[i]["status"])
		statusJ := asString(clubs[j]["status"])
		if statusI != statusJ {
			return statusI < statusJ
		}
		return asString(clubs[i]["name"]) < asString(clubs[j]["name"])
	})

	return clubs, nil
}

// AdminGetClubs returns all clubs including pending/rejected.
func AdminGetClubs(c *gin.Context) {
	clubs, err := readAllClubs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clubs"})
		return
	}
	c.JSON(http.StatusOK, clubs)
}

// AdminCreateClub allows admin to create an approved club directly.
func AdminCreateClub(c *gin.Context) {
	adminID, err := getRequesterUID(c)
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
	adminName := getUserDisplayNameByUID(c, adminID, nil)
	club := map[string]interface{}{
		"name":            input.Name,
		"description":     input.Description,
		"icon":            input.Icon,
		"color":           input.Color,
		"members":         []string{adminID},
		"created_by":      adminID,
		"created_by_name": adminName,
		"status":          "approved",
		"created_at":      now,
		"updated_at":      now,
		"moderated_at":    now,
		"moderated_by":    adminID,
	}

	docRef, _, err := repositories.FirestoreClient.Collection("clubs").Add(c.Request.Context(), club)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create club"})
		return
	}

	club["id"] = docRef.ID
	logAdminAction(c, adminID, "club.created", "club", docRef.ID, input.Name, nil)
	c.JSON(http.StatusCreated, club)
}

// AdminDeleteClub removes club.
func AdminDeleteClub(c *gin.Context) {
	id := c.Param("id")
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	targetName := id
	if doc, getErr := repositories.FirestoreClient.Collection("clubs").Doc(id).Get(c.Request.Context()); getErr == nil {
		name := asString(doc.Data()["name"])
		if name != "" {
			targetName = name
		}
	}

	_, err = repositories.FirestoreClient.Collection("clubs").Doc(id).Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete club"})
		return
	}

	logAdminAction(c, adminID, "club.deleted", "club", id, targetName, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Club deleted"})
}

// AdminUpdateClub edits club and allows changing moderation status.
func AdminUpdateClub(c *gin.Context) {
	id := c.Param("id")
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	oldStatus := ""
	oldName := id
	if doc, getErr := repositories.FirestoreClient.Collection("clubs").Doc(id).Get(c.Request.Context()); getErr == nil {
		oldStatus = getClubStatus(doc.Data())
		if name := asString(doc.Data()["name"]); name != "" {
			oldName = name
		}
	}

	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	clubInput := ClubInput{Name: input.Name, Description: input.Description, Icon: input.Icon, Color: input.Color}
	if !normalizeClubInput(&clubInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and description are required"})
		return
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "approved"
	}
	if status != "approved" && status != "pending" && status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	update := map[string]interface{}{
		"name":         clubInput.Name,
		"description":  clubInput.Description,
		"icon":         clubInput.Icon,
		"color":        clubInput.Color,
		"status":       status,
		"updated_at":   time.Now(),
		"moderated_at": time.Now(),
		"moderated_by": adminID,
	}

	_, err = repositories.FirestoreClient.Collection("clubs").Doc(id).Set(c.Request.Context(), update, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update club"})
		return
	}

	logAdminAction(c, adminID, "club.updated", "club", id, clubInput.Name, map[string]interface{}{
		"old_status": oldStatus,
		"new_status": status,
		"old_name":   oldName,
		"new_name":   clubInput.Name,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Club updated"})
}

// AdminGetClubRequests returns only pending club requests.
func AdminGetClubRequests(c *gin.Context) {
	clubs, err := readAllClubs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch club requests"})
		return
	}

	var pending []map[string]interface{}
	for _, club := range clubs {
		if asString(club["status"]) == "pending" {
			pending = append(pending, club)
		}
	}

	c.JSON(http.StatusOK, pending)
}

// AdminApproveClubRequest approves club request.
func AdminApproveClubRequest(c *gin.Context) {
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	targetName := id
	if doc, getErr := repositories.FirestoreClient.Collection("clubs").Doc(id).Get(c.Request.Context()); getErr == nil {
		name := asString(doc.Data()["name"])
		if name != "" {
			targetName = name
		}
	}

	_, err = repositories.FirestoreClient.Collection("clubs").Doc(id).Set(c.Request.Context(), map[string]interface{}{
		"status":             "approved",
		"moderated_by":       adminID,
		"moderated_by_name":  getUserDisplayNameByUID(c, adminID, nil),
		"moderated_at":       time.Now(),
		"moderation_comment": "",
		"updated_at":         time.Now(),
	}, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve club request"})
		return
	}

	logAdminAction(c, adminID, "club.request.approved", "club", id, targetName, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Club request approved"})
}

// AdminRejectClubRequest rejects club request.
func AdminRejectClubRequest(c *gin.Context) {
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	targetName := id
	if doc, getErr := repositories.FirestoreClient.Collection("clubs").Doc(id).Get(c.Request.Context()); getErr == nil {
		name := asString(doc.Data()["name"])
		if name != "" {
			targetName = name
		}
	}
	var input struct {
		Comment string `json:"comment"`
	}
	_ = c.ShouldBindJSON(&input)

	_, err = repositories.FirestoreClient.Collection("clubs").Doc(id).Set(c.Request.Context(), map[string]interface{}{
		"status":             "rejected",
		"moderated_by":       adminID,
		"moderated_by_name":  getUserDisplayNameByUID(c, adminID, nil),
		"moderated_at":       time.Now(),
		"moderation_comment": strings.TrimSpace(input.Comment),
		"updated_at":         time.Now(),
	}, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject club request"})
		return
	}

	logAdminAction(c, adminID, "club.request.rejected", "club", id, targetName, map[string]interface{}{
		"comment": strings.TrimSpace(input.Comment),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Club request rejected"})
}
