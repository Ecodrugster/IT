package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

type GradeInput struct {
	StudentID  string `json:"student_id" binding:"required"`
	ScheduleID string `json:"schedule_id" binding:"required"`
	Value      int    `json:"value" binding:"required"`
	Comment    string `json:"comment"`
	LessonDate string `json:"lesson_date"`
}

// AddGrade allows teachers/admins to add a grade bound to a schedule pair.
func AddGrade(c *gin.Context) {
	teacherID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, err := getUserRoleByUID(c, teacherID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to detect teacher role"})
		return
	}
	if !isTeacherLikeRole(role) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Teacher privileges required"})
		return
	}

	var input GradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.StudentID = strings.TrimSpace(input.StudentID)
	input.ScheduleID = strings.TrimSpace(input.ScheduleID)
	input.Comment = strings.TrimSpace(input.Comment)
	input.LessonDate = strings.TrimSpace(input.LessonDate)

	if input.StudentID == "" || input.ScheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id and schedule_id are required"})
		return
	}
	if input.Value < 1 || input.Value > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Grade value should be between 1 and 12"})
		return
	}

	studentDoc, err := repositories.FirestoreClient.Collection("users").Doc(input.StudentID).Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student not found"})
		return
	}
	studentRole := asString(studentDoc.Data()["role"])
	if studentRole != "" && studentRole != "student" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Selected user is not a student"})
		return
	}
	studentGroup := asString(studentDoc.Data()["group_name"])
	if studentGroup == "" {
		studentGroup = asString(studentDoc.Data()["group"])
	}

	scheduleDoc, err := repositories.FirestoreClient.Collection("schedule").Doc(input.ScheduleID).Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule item not found"})
		return
	}
	scheduleData := scheduleDoc.Data()
	scheduleTeacherID := asString(scheduleData["teacher_id"])
	if role == "teacher" && scheduleTeacherID != teacherID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can add grades only for your own pairs"})
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

	grade := map[string]interface{}{
		"student_id":   input.StudentID,
		"teacher_id":   teacherID,
		"teacher_name": teacherName,
		"schedule_id":  input.ScheduleID,
		"subject":      asString(scheduleData["subject"]),
		"group_name":   asString(scheduleData["group_name"]),
		"room":         asString(scheduleData["room"]),
		"day_of_week":  asInt(scheduleData["day_of_week"]),
		"pair_number":  asInt(scheduleData["pair_number"]),
		"starts_at":    asString(scheduleData["starts_at"]),
		"ends_at":      asString(scheduleData["ends_at"]),
		"value":        input.Value,
		"comment":      input.Comment,
		"lesson_date":  lessonDate,
		"created_at":   time.Now(),
	}

	_, _, err = repositories.FirestoreClient.Collection("grades").Add(c.Request.Context(), grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add grade"})
		return
	}

	// Send notification to student
	notifTitle := "Получена новая оценка"
	notifMsg := fmt.Sprintf("Вы получили оценку %d по предмету %s от преподавателя %s", input.Value, asString(grade["subject"]), teacherName)
	if input.Comment != "" {
		notifMsg += fmt.Sprintf(" (%s)", input.Comment)
	}
	notifLink := "/profile/grades"
	_ = CreateNotification(c.Request.Context(), input.StudentID, "grade", notifTitle, notifMsg, notifLink)

	c.JSON(http.StatusCreated, gin.H{"message": "Grade added successfully"})
}

// GetUserGrades returns student's grades, and allows teachers/admins to query by student_id.
func GetUserGrades(c *gin.Context) {
	requesterID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	studentID := strings.TrimSpace(c.Query("student_id"))
	if studentID == "" {
		studentID = requesterID
	}

	if studentID != requesterID {
		role, err := getUserRoleByUID(c, requesterID)
		if err != nil || !isTeacherLikeRole(role) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot view grades of other students"})
			return
		}

		if role == "teacher" {
			studentDoc, studentErr := repositories.FirestoreClient.Collection("users").Doc(studentID).Get(c.Request.Context())
			if studentErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Student not found"})
				return
			}
			studentGroup := asString(studentDoc.Data()["group_name"])
			if studentGroup == "" {
				studentGroup = asString(studentDoc.Data()["group"])
			}
			if studentGroup == "" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Student has no assigned group"})
				return
			}

			canAccess := false
			scheduleIter := repositories.FirestoreClient.Collection("schedule").Where("teacher_id", "==", requesterID).Documents(c.Request.Context())
			for {
				doc, iterErr := scheduleIter.Next()
				if iterErr == iterator.Done {
					break
				}
				if iterErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate teacher access"})
					return
				}
				groupName := asString(doc.Data()["group_name"])
				if strings.EqualFold(groupName, studentGroup) {
					canAccess = true
					break
				}
			}

			if !canAccess {
				c.JSON(http.StatusForbidden, gin.H{"error": "You cannot view grades of students outside your groups"})
				return
			}
		}
	}

	iter := repositories.FirestoreClient.Collection("grades").Where("student_id", "==", studentID).OrderBy("created_at", repositories.Descending).Documents(c.Request.Context())

	var grades []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grades"})
			return
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		grades = append(grades, data)
	}

	c.JSON(http.StatusOK, grades)
}
