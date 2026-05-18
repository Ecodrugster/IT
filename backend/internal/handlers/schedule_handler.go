package handlers

import (
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

type ScheduleInput struct {
	Subject    string `json:"subject" binding:"required"`
	GroupName  string `json:"group_name" binding:"required"`
	TeacherID  string `json:"teacher_id" binding:"required"`
	DayOfWeek  int    `json:"day_of_week" binding:"required"`
	PairNumber int    `json:"pair_number" binding:"required"`
	StartsAt   string `json:"starts_at" binding:"required"`
	EndsAt     string `json:"ends_at"`
	Room       string `json:"room"`
}

func validateScheduleInput(input *ScheduleInput) bool {
	input.Subject = strings.TrimSpace(input.Subject)
	input.GroupName = strings.TrimSpace(input.GroupName)
	input.TeacherID = strings.TrimSpace(input.TeacherID)
	input.StartsAt = strings.TrimSpace(input.StartsAt)
	input.EndsAt = strings.TrimSpace(input.EndsAt)
	input.Room = strings.TrimSpace(input.Room)

	if input.Subject == "" || input.GroupName == "" || input.TeacherID == "" {
		return false
	}

	if input.DayOfWeek < 1 || input.DayOfWeek > 7 {
		return false
	}

	if input.PairNumber < 1 || input.PairNumber > 10 {
		return false
	}

	if _, err := time.Parse("15:04", input.StartsAt); err != nil {
		return false
	}

	if input.EndsAt != "" {
		if _, err := time.Parse("15:04", input.EndsAt); err != nil {
			return false
		}
	}

	return true
}

func getTeacherMeta(c *gin.Context, teacherID string) (string, string, error) {
	doc, err := repositories.FirestoreClient.Collection("users").Doc(teacherID).Get(c.Request.Context())
	if err != nil {
		return "", "", err
	}

	data := doc.Data()
	role := asString(data["role"])
	if role == "" {
		role = "student"
	}
	if !isTeacherLikeRole(role) {
		return "", "", errors.New("invalid teacher role")
	}

	teacherName := asString(data["display_name"])
	if teacherName == "" {
		teacherName = asString(data["displayName"])
	}
	if teacherName == "" {
		teacherName = asString(data["email"])
	}

	return teacherName, role, nil
}

func getUserGroup(c *gin.Context, uid string) string {
	doc, err := repositories.FirestoreClient.Collection("users").Doc(uid).Get(c.Request.Context())
	if err != nil {
		return ""
	}

	data := doc.Data()
	group := asString(data["group"])
	if group == "" {
		group = asString(data["group_name"])
	}
	return group
}

func hasScheduleConflict(c *gin.Context, input ScheduleInput, skipID string) (bool, string, error) {
	iter := repositories.FirestoreClient.Collection("schedule").Documents(c.Request.Context())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, "", err
		}

		if skipID != "" && doc.Ref.ID == skipID {
			continue
		}

		data := doc.Data()
		day := asInt(data["day_of_week"])
		pair := asInt(data["pair_number"])
		if day != input.DayOfWeek || pair != input.PairNumber {
			continue
		}

		groupName := asString(data["group_name"])
		teacherID := asString(data["teacher_id"])

		if strings.EqualFold(groupName, input.GroupName) {
			return true, "A pair already exists for this group at the same time", nil
		}
		if teacherID == input.TeacherID {
			return true, "Teacher already has another pair at this time", nil
		}
	}

	return false, "", nil
}

func readSchedule(c *gin.Context) ([]map[string]interface{}, error) {
	iter := repositories.FirestoreClient.Collection("schedule").Documents(c.Request.Context())
	var items []map[string]interface{}

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
		items = append(items, data)
	}

	sort.Slice(items, func(i, j int) bool {
		dayI := asInt(items[i]["day_of_week"])
		dayJ := asInt(items[j]["day_of_week"])
		if dayI != dayJ {
			return dayI < dayJ
		}
		pairI := asInt(items[i]["pair_number"])
		pairJ := asInt(items[j]["pair_number"])
		if pairI != pairJ {
			return pairI < pairJ
		}
		return asString(items[i]["subject"]) < asString(items[j]["subject"])
	})

	return items, nil
}

func GetSchedule(c *gin.Context) {
	uid, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, err := getUserRoleByUID(c, uid)
	if err != nil {
		role = "student"
	}

	groupFilter := strings.TrimSpace(c.Query("group"))
	teacherFilter := strings.TrimSpace(c.Query("teacher_id"))

	if role == "student" {
		ownGroup := getUserGroup(c, uid)
		if ownGroup == "" {
			c.JSON(http.StatusOK, []map[string]interface{}{})
			return
		}
		groupFilter = ownGroup
	} else if groupFilter == "" && !isAdminLikeRole(role) {
		groupFilter = getUserGroup(c, uid)
	}

	items, err := readSchedule(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}

	var filtered []map[string]interface{}
	for _, item := range items {
		if groupFilter != "" && !strings.EqualFold(asString(item["group_name"]), groupFilter) {
			continue
		}
		if teacherFilter != "" && asString(item["teacher_id"]) != teacherFilter {
			continue
		}
		filtered = append(filtered, item)
	}

	c.JSON(http.StatusOK, filtered)
}

func GetTeacherSchedule(c *gin.Context) {
	uid, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	role, _ := getUserRoleByUID(c, uid)

	items, err := readSchedule(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teacher schedule"})
		return
	}

	if role == "admin" {
		c.JSON(http.StatusOK, items)
		return
	}

	var mine []map[string]interface{}
	for _, item := range items {
		if asString(item["teacher_id"]) == uid {
			mine = append(mine, item)
		}
	}

	c.JSON(http.StatusOK, mine)
}

func AdminGetSchedule(c *gin.Context) {
	items, err := readSchedule(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func AdminCreateSchedule(c *gin.Context) {
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !validateScheduleInput(&input) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule payload"})
		return
	}

	teacherName, _, err := getTeacherMeta(c, input.TeacherID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Teacher not found or has invalid role"})
		return
	}

	conflict, reason, err := hasScheduleConflict(c, input, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate schedule"})
		return
	}
	if conflict {
		c.JSON(http.StatusConflict, gin.H{"error": reason})
		return
	}

	now := time.Now()
	item := map[string]interface{}{
		"subject":      input.Subject,
		"group_name":   input.GroupName,
		"teacher_id":   input.TeacherID,
		"teacher_name": teacherName,
		"day_of_week":  input.DayOfWeek,
		"pair_number":  input.PairNumber,
		"starts_at":    input.StartsAt,
		"ends_at":      input.EndsAt,
		"room":         input.Room,
		"created_by":   adminID,
		"created_at":   now,
		"updated_at":   now,
	}

	docRef, _, err := repositories.FirestoreClient.Collection("schedule").Add(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule item"})
		return
	}

	item["id"] = docRef.ID
	logAdminAction(c, adminID, "schedule.created", "schedule", docRef.ID, input.Subject, map[string]interface{}{
		"group_name":   input.GroupName,
		"teacher_id":   input.TeacherID,
		"teacher_name": teacherName,
		"day_of_week":  input.DayOfWeek,
		"pair_number":  input.PairNumber,
	})
	c.JSON(http.StatusCreated, item)
}

func AdminUpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule id is required"})
		return
	}
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !validateScheduleInput(&input) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule payload"})
		return
	}

	teacherName, _, err := getTeacherMeta(c, input.TeacherID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Teacher not found or has invalid role"})
		return
	}

	conflict, reason, err := hasScheduleConflict(c, input, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate schedule"})
		return
	}
	if conflict {
		c.JSON(http.StatusConflict, gin.H{"error": reason})
		return
	}

	ref := repositories.FirestoreClient.Collection("schedule").Doc(id)
	existingDoc, err := ref.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule item not found"})
		return
	}
	oldData := existingDoc.Data()

	update := map[string]interface{}{
		"subject":      input.Subject,
		"group_name":   input.GroupName,
		"teacher_id":   input.TeacherID,
		"teacher_name": teacherName,
		"day_of_week":  input.DayOfWeek,
		"pair_number":  input.PairNumber,
		"starts_at":    input.StartsAt,
		"ends_at":      input.EndsAt,
		"room":         input.Room,
		"updated_at":   time.Now(),
	}

	_, err = ref.Set(c.Request.Context(), update, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule item"})
		return
	}

	logAdminAction(c, adminID, "schedule.updated", "schedule", id, input.Subject, map[string]interface{}{
		"old_subject":  oldData["subject"],
		"new_subject":  input.Subject,
		"group_name":   input.GroupName,
		"teacher_id":   input.TeacherID,
		"teacher_name": teacherName,
		"day_of_week":  input.DayOfWeek,
		"pair_number":  input.PairNumber,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Schedule item updated"})
}

func AdminDeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule id is required"})
		return
	}
	adminID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	targetName := id
	details := map[string]interface{}{}
	if doc, getErr := repositories.FirestoreClient.Collection("schedule").Doc(id).Get(c.Request.Context()); getErr == nil {
		item := doc.Data()
		if subject := asString(item["subject"]); subject != "" {
			targetName = subject
		}
		details["group_name"] = asString(item["group_name"])
		details["teacher_id"] = asString(item["teacher_id"])
		details["day_of_week"] = asInt(item["day_of_week"])
		details["pair_number"] = asInt(item["pair_number"])
	}

	_, err = repositories.FirestoreClient.Collection("schedule").Doc(id).Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule item"})
		return
	}

	logAdminAction(c, adminID, "schedule.deleted", "schedule", id, targetName, details)

	c.JSON(http.StatusOK, gin.H{"message": "Schedule item deleted"})
}
