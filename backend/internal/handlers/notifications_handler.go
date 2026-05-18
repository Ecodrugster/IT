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
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Notification represents a system notification stored in Firestore
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"` // "message" | "grade" | "attendance" | "news" | "system"
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Link      string    `json:"link"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateNotification is a helper to easily write a new notification to Firestore
func CreateNotification(ctx context.Context, userID, notifType, title, message, link string) error {
	userID = strings.TrimSpace(userID)
	notifType = strings.TrimSpace(notifType)
	title = strings.TrimSpace(title)
	message = strings.TrimSpace(message)
	link = strings.TrimSpace(link)

	if userID == "" {
		return fmt.Errorf("user_id is required")
	}
	if notifType == "" {
		notifType = "system"
	}

	entry := map[string]interface{}{
		"user_id":    userID,
		"type":       notifType,
		"title":      title,
		"message":    message,
		"link":       link,
		"read":       false,
		"created_at": time.Now(),
	}

	_, _, err := repositories.FirestoreClient.Collection("notifications").Add(ctx, entry)
	return err
}

// GetNotifications returns all notifications for the authenticated user, sorted in Go by date descending
func GetNotifications(c *gin.Context) {
	userID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	iter := repositories.FirestoreClient.Collection("notifications").
		Where("user_id", "==", userID).
		Documents(c.Request.Context())

	var list []Notification
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
			return
		}

		data := doc.Data()
		
		var createdAt time.Time
		if t, ok := data["created_at"].(time.Time); ok {
			createdAt = t
		}

		notif := Notification{
			ID:        doc.Ref.ID,
			UserID:    asString(data["user_id"]),
			Type:      asString(data["type"]),
			Title:     asString(data["title"]),
			Message:   asString(data["message"]),
			Link:      asString(data["link"]),
			Read:      data["read"] == true,
			CreatedAt: createdAt,
		}
		list = append(list, notif)
	}

	// Sort by CreatedAt descending in Go to avoid Firestore Index requirements
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.After(list[j].CreatedAt)
	})

	// Limit to recent 100 notifications
	if len(list) > 100 {
		list = list[:100]
	}

	c.JSON(http.StatusOK, list)
}

// MarkNotificationAsRead marks a single notification as read
func MarkNotificationAsRead(c *gin.Context) {
	userID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notifID := strings.TrimSpace(c.Param("id"))
	if notifID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID is required"})
		return
	}

	docRef := repositories.FirestoreClient.Collection("notifications").Doc(notifID)
	doc, err := docRef.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	if asString(doc.Data()["user_id"]) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	_, err = docRef.Update(c.Request.Context(), []firestore.Update{
		{Path: "read", Value: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAllNotificationsAsRead marks all notifications for user as read using Firestore Batch
func MarkAllNotificationsAsRead(c *gin.Context) {
	userID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	iter := repositories.FirestoreClient.Collection("notifications").
		Where("user_id", "==", userID).
		Where("read", "==", false).
		Documents(c.Request.Context())

	batch := repositories.FirestoreClient.Batch()
	count := 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query unread notifications"})
			return
		}

		batch.Set(doc.Ref, map[string]interface{}{"read": true}, firestore.MergeAll)
		count++

		// Firestore batch limit is 500 operations
		if count >= 500 {
			_, err = batch.Commit(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write batch updates"})
				return
			}
			batch = repositories.FirestoreClient.Batch()
			count = 0
		}
	}

	if count > 0 {
		_, err = batch.Commit(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write batch updates"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

// GetUnifiedUnreadCount calculates unread firestore notifications, unread messages, and combined count
func GetUnifiedUnreadCount(c *gin.Context) {
	userID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 1. Calculate Firestore unread count
	iter := repositories.FirestoreClient.Collection("notifications").
		Where("user_id", "==", userID).
		Where("read", "==", false).
		Documents(c.Request.Context())

	notifCount := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}
		notifCount++
	}

	// 2. Calculate Mongo chat unread count
	var chatCount int64 = 0
	if repositories.IsMongoChatReady() {
		filter := bson.M{
			"receiver_id": userID,
			"read":        false,
		}
		var err error
		chatCount, err = repositories.ChatMessagesCollection.CountDocuments(c.Request.Context(), filter)
		if err != nil {
			chatCount = 0
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_unread":  notifCount + int(chatCount),
		"notifications": notifCount,
		"chats":         int(chatCount),
	})
}
