package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

type AttendanceInput struct {
	StudentID  string `json:"student_id" binding:"required"`
	ScheduleID string `json:"schedule_id" binding:"required"`
	LessonDate string `json:"lesson_date"`
	Status     string `json:"status" binding:"required"` // present | absent | late | excused
	Comment    string `json:"comment"`
}

func normalizeAttendanceStatus(raw string) string {
	status := strings.ToLower(strings.TrimSpace(raw))
	switch status {
	case "present", "absent", "late", "excused":
		return status
	default:
		return ""
	}
}

func MarkAttendance(c *gin.Context) {
	teacherID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, err := getUserRoleByUID(c, teacherID)
	if err != nil || !isTeacherLikeRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Teacher privileges required"})
		return
	}

	var input AttendanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.StudentID = strings.TrimSpace(input.StudentID)
	input.ScheduleID = strings.TrimSpace(input.ScheduleID)
	input.LessonDate = strings.TrimSpace(input.LessonDate)
	input.Status = normalizeAttendanceStatus(input.Status)
	input.Comment = strings.TrimSpace(input.Comment)

	if input.StudentID == "" || input.ScheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id and schedule_id are required"})
		return
	}
	if input.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be one of: present, absent, late, excused"})
		return
	}

	studentDoc, err := repositories.FirestoreClient.Collection("users").Doc(input.StudentID).Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student not found"})
		return
	}
	studentData := studentDoc.Data()
	studentRole := normalizeRole(asString(studentData["role"]))
	if studentRole != "student" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Selected user is not a student"})
		return
	}
	studentGroup := asString(studentData["group_name"])
	if studentGroup == "" {
		studentGroup = asString(studentData["group"])
	}

	scheduleDoc, err := repositories.FirestoreClient.Collection("schedule").Doc(input.ScheduleID).Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule item not found"})
		return
	}
	scheduleData := scheduleDoc.Data()
	scheduleTeacherID := asString(scheduleData["teacher_id"])
	if role == "teacher" && scheduleTeacherID != teacherID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can mark attendance only for your own pairs"})
		return
	}

	scheduleGroup := asString(scheduleData["group_name"])
	if studentGroup != "" && scheduleGroup != "" && !strings.EqualFold(studentGroup, scheduleGroup) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student is not assigned to the pair group"})
		return
	}

	lessonDate := input.LessonDate
	if lessonDate == "" {
		lessonDate = time.Now().Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", lessonDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lesson_date must be in YYYY-MM-DD format"})
		return
	}

	teacherName := asString(scheduleData["teacher_name"])
	if teacherName == "" {
		teacherName = getUserDisplayNameByUID(c, teacherID, nil)
	}

	studentName := getDisplayNameFromUserData(studentData, input.StudentID)
	now := time.Now()
	entry := map[string]interface{}{
		"student_id":   input.StudentID,
		"student_name": studentName,
		"teacher_id":   teacherID,
		"teacher_name": teacherName,
		"schedule_id":  input.ScheduleID,
		"subject":      asString(scheduleData["subject"]),
		"group_name":   scheduleGroup,
		"room":         asString(scheduleData["room"]),
		"day_of_week":  asInt(scheduleData["day_of_week"]),
		"pair_number":  asInt(scheduleData["pair_number"]),
		"starts_at":    asString(scheduleData["starts_at"]),
		"ends_at":      asString(scheduleData["ends_at"]),
		"lesson_date":  lessonDate,
		"status":       input.Status,
		"comment":      input.Comment,
		"updated_at":   now,
	}

	var existingDocID string
	var oldStatus string
	iter := repositories.FirestoreClient.Collection("attendance").Where("student_id", "==", input.StudentID).Documents(c.Request.Context())
	for {
		doc, iterErr := iter.Next()
		if iterErr == iterator.Done {
			break
		}
		if iterErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing attendance"})
			return
		}
		data := doc.Data()
		if asString(data["schedule_id"]) == input.ScheduleID && asString(data["lesson_date"]) == lessonDate {
			existingDocID = doc.Ref.ID
			oldStatus = asString(data["status"])
			break
		}
	}

	// Helper to send notification
	sendAttendanceNotif := func() {
		statusRU := map[string]string{
			"absent":  "пропуск (отсутствие)",
			"late":    "опоздание",
			"excused": "пропуск по уважительной причине",
			"present": "присутствие",
		}[input.Status]
		if statusRU == "" {
			statusRU = input.Status
		}

		notifTitle := "Изменение статуса посещаемости"
		notifMsg := fmt.Sprintf("Отмечено: %s по предмету %s за %s. Преподаватель: %s", statusRU, asString(entry["subject"]), lessonDate, teacherName)
		if input.Comment != "" {
			notifMsg += fmt.Sprintf(" (%s)", input.Comment)
		}
		notifLink := "/schedule"
		_ = CreateNotification(c.Request.Context(), input.StudentID, "attendance", notifTitle, notifMsg, notifLink)
	}

	if existingDocID != "" {
		_, err = repositories.FirestoreClient.Collection("attendance").Doc(existingDocID).Set(c.Request.Context(), entry, repositories.MergeAll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
			return
		}
		UpdateStudentGamification(c.Request.Context(), input.StudentID, oldStatus, input.Status)
		sendAttendanceNotif()
		c.JSON(http.StatusOK, gin.H{"message": "Attendance updated", "id": existingDocID})
		return
	}

	entry["created_at"] = now
	docRef, _, err := repositories.FirestoreClient.Collection("attendance").Add(c.Request.Context(), entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark attendance"})
		return
	}

	UpdateStudentGamification(c.Request.Context(), input.StudentID, "", input.Status)
	sendAttendanceNotif()
	c.JSON(http.StatusCreated, gin.H{"message": "Attendance marked", "id": docRef.ID})
}

func GetAttendance(c *gin.Context) {
	requesterID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, _ := getUserRoleByUID(c, requesterID)
	role = normalizeRole(role)

	studentID := strings.TrimSpace(c.Query("student_id"))
	scheduleID := strings.TrimSpace(c.Query("schedule_id"))
	lessonDate := strings.TrimSpace(c.Query("lesson_date"))
	groupName := strings.TrimSpace(c.Query("group"))

	if studentID == "" && role == "student" {
		studentID = requesterID
	}
	if studentID != "" && studentID != requesterID && role == "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot view attendance of other users"})
		return
	}

	var iter *firestore.DocumentIterator
	if studentID != "" {
		iter = repositories.FirestoreClient.Collection("attendance").Where("student_id", "==", studentID).Documents(c.Request.Context())
	} else {
		iter = repositories.FirestoreClient.Collection("attendance").Documents(c.Request.Context())
	}

	var items []map[string]interface{}
	for {
		doc, iterErr := iter.Next()
		if iterErr == iterator.Done {
			break
		}
		if iterErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance"})
			return
		}

		data := doc.Data()
		if scheduleID != "" && asString(data["schedule_id"]) != scheduleID {
			continue
		}
		if lessonDate != "" && asString(data["lesson_date"]) != lessonDate {
			continue
		}
		if groupName != "" && !strings.EqualFold(asString(data["group_name"]), groupName) {
			continue
		}

		// Teacher can only view rows marked by them.
		if role == "teacher" && asString(data["teacher_id"]) != requesterID {
			continue
		}

		data["id"] = doc.Ref.ID
		items = append(items, data)
	}

	sort.Slice(items, func(i, j int) bool {
		dateI := asString(items[i]["lesson_date"])
		dateJ := asString(items[j]["lesson_date"])
		if dateI != dateJ {
			return dateI > dateJ
		}
		pairI := asInt(items[i]["pair_number"])
		pairJ := asInt(items[j]["pair_number"])
		if pairI != pairJ {
			return pairI < pairJ
		}
		ti, okI := items[i]["updated_at"].(time.Time)
		tj, okJ := items[j]["updated_at"].(time.Time)
		if !okI || !okJ {
			return false
		}
		return ti.After(tj)
	})

	c.JSON(http.StatusOK, items)
}

// UpdateStudentGamification updates a student's coins, stars, and streak based on attendance mark changes.
func UpdateStudentGamification(ctx context.Context, studentID, oldStatus, newStatus string) {
	studentID = strings.TrimSpace(studentID)
	if studentID == "" {
		return
	}

	studentRef := repositories.FirestoreClient.Collection("users").Doc(studentID)
	studentDoc, err := studentRef.Get(ctx)
	if err != nil {
		return
	}

	data := studentDoc.Data()

	extractInt := func(val interface{}) int {
		switch v := val.(type) {
		case int64:
			return int(v)
		case float64:
			return int(v)
		case int:
			return v
		case int32:
			return int(v)
		case float32:
			return int(v)
		default:
			return 0
		}
	}

	coins := extractInt(data["coins"])
	stars := extractInt(data["stars"])
	streak := extractInt(data["streak"])

	isActive := func(status string) bool {
		return status == "present" || status == "late"
	}
	isInactive := func(status string) bool {
		return status == "absent"
	}

	updated := false
	gotStar := false
	gotCoins := 0

	if oldStatus == "" {
		if isActive(newStatus) {
			coins += 5
			gotCoins = 5
			streak++
			if streak >= 5 {
				stars++
				gotStar = true
				streak = 0
			}
			updated = true
		} else if isInactive(newStatus) {
			streak = 0
			updated = true
		}
	} else {
		if oldStatus != newStatus {
			if isActive(oldStatus) && !isActive(newStatus) {
				coins -= 5
				if coins < 0 {
					coins = 0
				}
				streak = 0
				updated = true
			} else if !isActive(oldStatus) && isActive(newStatus) {
				coins += 5
				gotCoins = 5
				streak++
				if streak >= 5 {
					stars++
					gotStar = true
					streak = 0
				}
				updated = true
			} else if isInactive(oldStatus) && !isInactive(newStatus) && !isActive(newStatus) {
				// Changed from absent to excused, streak resets or stays 0
				streak = 0
				updated = true
			}
		}
	}

	if updated {
		_, _ = studentRef.Set(ctx, map[string]interface{}{
			"coins":  coins,
			"stars":  stars,
			"streak": streak,
		}, firestore.MergeAll)

		if gotCoins > 0 {
			notifTitle := "Получены монеты! 🪙"
			notifMsg := fmt.Sprintf("Вы получили %d монет(ы) за посещение занятия. Текущий баланс: %d монет(ы)", gotCoins, coins)
			_ = CreateNotification(ctx, studentID, "system", notifTitle, notifMsg, "/store")
		}
		if gotStar {
			notifTitle := "Получена звезда! 🌟"
			notifMsg := fmt.Sprintf("Поздравляем! 5 посещений занятий подряд без пропусков! Вы заработали 1 звезду. Текущий баланс: %d звезд(ы)", stars)
			_ = CreateNotification(ctx, studentID, "system", notifTitle, notifMsg, "/store")
		}
	}
}
