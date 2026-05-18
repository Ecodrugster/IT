package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
)

func getRequesterUID(c *gin.Context) (string, error) {
	uid := strings.TrimSpace(c.GetString("firebase_uid"))
	if uid == "" {
		return "", fmt.Errorf("missing authenticated user id")
	}
	return uid, nil
}

func getUserRoleByUID(c *gin.Context, uid string) (string, error) {
	doc, err := repositories.FirestoreClient.Collection("users").Doc(uid).Get(c.Request.Context())
	if err != nil {
		return "", err
	}

	role, _ := doc.Data()["role"].(string)
	if role == "" {
		return "student", nil
	}

	return role, nil
}

func isAdminLikeRole(role string) bool {
	return role == "admin"
}

func isTeacherLikeRole(role string) bool {
	return role == "teacher" || role == "admin"
}

func asString(v interface{}) string {
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

func asInt(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float32:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}
