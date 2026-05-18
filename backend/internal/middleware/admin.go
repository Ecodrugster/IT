package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
)

// AdminRequired allows only admin users.
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		firebaseUID, exists := c.Get("firebase_uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		doc, err := repositories.FirestoreClient.Collection("users").Doc(firebaseUID.(string)).Get(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. User record not found."})
			c.Abort()
			return
		}

		role, _ := doc.Data()["role"].(string)
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required."})
			c.Abort()
			return
		}

		c.Next()
	}
}

// TeacherRequired allows teachers and admins.
func TeacherRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		firebaseUID, exists := c.Get("firebase_uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		doc, err := repositories.FirestoreClient.Collection("users").Doc(firebaseUID.(string)).Get(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied."})
			c.Abort()
			return
		}

		role, _ := doc.Data()["role"].(string)
		if role != "teacher" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Teacher privileges required."})
			c.Abort()
			return
		}

		c.Next()
	}
}
