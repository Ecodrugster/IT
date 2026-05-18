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

func getDisplayNameFromUserData(data map[string]interface{}, fallback string) string {
	displayName := asString(data["display_name"])
	if displayName == "" {
		displayName = asString(data["displayName"])
	}
	if displayName == "" {
		displayName = asString(data["name"])
	}
	if displayName == "" {
		displayName = asString(data["email"])
	}
	if displayName == "" {
		displayName = fallback
	}
	return displayName
}

func getUserDisplayNameByUID(c *gin.Context, uid string, cache map[string]string) string {
	uid = strings.TrimSpace(uid)
	if uid == "" {
		return ""
	}

	if cache != nil {
		if name, ok := cache[uid]; ok {
			return name
		}
	}

	doc, err := repositories.FirestoreClient.Collection("users").Doc(uid).Get(c.Request.Context())
	if err != nil {
		if cache != nil {
			cache[uid] = uid
		}
		return uid
	}

	name := getDisplayNameFromUserData(doc.Data(), uid)
	if cache != nil {
		cache[uid] = name
	}
	return name
}

func logAdminAction(
	c *gin.Context,
	actorUID string,
	action string,
	targetType string,
	targetID string,
	targetName string,
	details map[string]interface{},
) {
	actorUID = strings.TrimSpace(actorUID)
	action = strings.TrimSpace(action)
	targetType = strings.TrimSpace(targetType)
	targetID = strings.TrimSpace(targetID)
	targetName = strings.TrimSpace(targetName)

	if actorUID == "" || action == "" {
		return
	}

	entry := map[string]interface{}{
		"actor_uid":   actorUID,
		"actor_name":  getUserDisplayNameByUID(c, actorUID, nil),
		"action":      action,
		"target_type": targetType,
		"target_id":   targetID,
		"target_name": targetName,
		"details":     details,
		"created_at":  time.Now(),
	}

	_, _, _ = repositories.FirestoreClient.Collection("admin_logs").Add(c.Request.Context(), entry)
}

func AdminGetDashboard(c *gin.Context) {
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)

	userStats := map[string]int{
		"total_users":    0,
		"total_students": 0,
		"total_teachers": 0,
		"total_admins":   0,
	}
	postStats := map[string]int{
		"total_posts":       0,
		"posts_last_7_days": 0,
	}
	clubStats := map[string]int{
		"active_clubs":          0,
		"pending_club_requests": 0,
	}
	newsStats := map[string]int{
		"total_news": 0,
	}

	usersIter := repositories.FirestoreClient.Collection("users").Documents(c.Request.Context())
	for {
		doc, err := usersIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read users for dashboard"})
			return
		}

		userStats["total_users"]++
		role := asString(doc.Data()["role"])
		if role == "" {
			role = "student"
		}

		switch role {
		case "admin":
			userStats["total_admins"]++
		case "teacher":
			userStats["total_teachers"]++
		default:
			userStats["total_students"]++
		}
	}

	postsIter := repositories.FirestoreClient.Collection("posts").Documents(c.Request.Context())
	for {
		doc, err := postsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read posts for dashboard"})
			return
		}

		postStats["total_posts"]++
		if createdAt, ok := doc.Data()["created_at"].(time.Time); ok && createdAt.After(weekAgo) {
			postStats["posts_last_7_days"]++
		}
	}

	clubsIter := repositories.FirestoreClient.Collection("clubs").Documents(c.Request.Context())
	for {
		doc, err := clubsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read clubs for dashboard"})
			return
		}

		status := getClubStatus(doc.Data())
		if status == "approved" {
			clubStats["active_clubs"]++
		}
		if status == "pending" {
			clubStats["pending_club_requests"]++
		}
	}

	newsIter := repositories.FirestoreClient.Collection("news").Documents(c.Request.Context())
	for {
		_, err := newsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read news for dashboard"})
			return
		}
		newsStats["total_news"]++
	}

	actionsIter := repositories.FirestoreClient.
		Collection("admin_logs").
		OrderBy("created_at", firestore.Desc).
		Limit(20).
		Documents(c.Request.Context())

	var actions []map[string]interface{}
	nameCache := map[string]string{}

	for {
		doc, err := actionsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read admin activity"})
			return
		}

		item := doc.Data()
		item["id"] = doc.Ref.ID

		actorUID := asString(item["actor_uid"])
		if actorUID != "" && asString(item["actor_name"]) == "" {
			item["actor_name"] = getUserDisplayNameByUID(c, actorUID, nameCache)
		}

		actions = append(actions, item)
	}

	sort.Slice(actions, func(i, j int) bool {
		ti, okI := actions[i]["created_at"].(time.Time)
		tj, okJ := actions[j]["created_at"].(time.Time)
		if !okI && !okJ {
			return false
		}
		if !okI {
			return false
		}
		if !okJ {
			return true
		}
		return ti.After(tj)
	})

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"users": userStats,
			"posts": postStats,
			"clubs": clubStats,
			"news":  newsStats,
		},
		"recent_actions": actions,
		"generated_at":   now,
	})
}
