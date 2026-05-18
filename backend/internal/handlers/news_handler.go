package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

type NewsItem struct {
	ID          string    `json:"id" firestore:"-"`
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	Category    string    `json:"category" firestore:"category"` // news, announcement, event, deadline
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
}

func CreateNews(c *gin.Context) {
	var item NewsItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.CreatedAt = time.Now()

	ref, _, err := repositories.FirestoreClient.Collection("news").Add(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news"})
		return
	}
	item.ID = ref.ID

	authorUID := strings.TrimSpace(c.GetString("firebase_uid"))
	if authorUID != "" {
		logAdminAction(c, authorUID, "news.created", "news", ref.ID, item.Title, map[string]interface{}{
			"category": item.Category,
		})
	}

	// Send notification to all users in background
	go func(newsItem NewsItem) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		iter := repositories.FirestoreClient.Collection("users").Documents(bgCtx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				break
			}

			userID := doc.Ref.ID
			notifTitle := "Новая публикация"
			switch newsItem.Category {
			case "announcement":
				notifTitle = "Важное объявление"
			case "event":
				notifTitle = "Новое мероприятие"
			case "deadline":
				notifTitle = "Дедлайн задания"
			}

			_ = CreateNotification(bgCtx, userID, "news", notifTitle, newsItem.Title, "/")
		}
	}(item)

	c.JSON(http.StatusCreated, item)
}

func GetNews(c *gin.Context) {
	category := c.Query("category")
	query := repositories.FirestoreClient.Collection("news").OrderBy("created_at", firestore.Desc)

	if category != "" {
		query = query.Where("category", "==", category)
	}

	iter := query.Documents(c.Request.Context())

	var news []NewsItem
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
			return
		}
		var item NewsItem
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		news = append(news, item)
	}

	c.JSON(http.StatusOK, news)
}
