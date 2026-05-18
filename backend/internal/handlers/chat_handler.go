package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/api/iterator"
)

const (
	defaultChatLimit      = 100
	maxChatLimit          = 300
	maxMessageLength      = 10000
	maxChatImageSizeBytes = 10 * 1024 * 1024

	defaultChatUploadRoot = "uploads"
	chatUploadSubDir      = "chat"

	messageTypeText  = "text"
	messageTypeImage = "image"
)

type chatMessageDocument struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	ChatID      string        `bson:"chat_id"`
	SenderID    string        `bson:"sender_id"`
	ReceiverID  string        `bson:"receiver_id"`
	MessageType string        `bson:"message_type,omitempty"`
	Text        string        `bson:"text,omitempty"`
	ImageURL    string        `bson:"image_url,omitempty"`
	Read        bool          `bson:"read"`
	CreatedAt   time.Time     `bson:"created_at"`
}

func ensureMongoChatReady(c *gin.Context) bool {
	if repositories.IsMongoChatReady() {
		return true
	}

	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": "Chat service is unavailable. MONGO_URL is not configured or MongoDB is unreachable.",
	})
	return false
}

func buildChatID(a, b string) string {
	users := []string{strings.TrimSpace(a), strings.TrimSpace(b)}
	sort.Strings(users)
	return users[0] + "_" + users[1]
}

func parseChatPeerID(c *gin.Context) string {
	userID := strings.TrimSpace(c.Query("user_id"))
	if userID == "" {
		userID = strings.TrimSpace(c.Query("uid"))
	}
	return userID
}

func parseChatReceiverID(input map[string]interface{}) string {
	receiverID := strings.TrimSpace(asString(input["receiver_id"]))
	if receiverID == "" {
		receiverID = strings.TrimSpace(asString(input["receiverId"]))
	}
	return receiverID
}

func parseChatReceiverIDFromForm(c *gin.Context) string {
	receiverID := strings.TrimSpace(c.PostForm("receiver_id"))
	if receiverID == "" {
		receiverID = strings.TrimSpace(c.PostForm("receiverId"))
	}
	return receiverID
}

func resolveMessageType(doc chatMessageDocument) string {
	if doc.MessageType != "" {
		return doc.MessageType
	}
	if strings.TrimSpace(doc.ImageURL) != "" {
		return messageTypeImage
	}
	return messageTypeText
}

func toChatMessageResponse(doc chatMessageDocument) gin.H {
	messageType := resolveMessageType(doc)
	imageURL := strings.TrimSpace(doc.ImageURL)

	return gin.H{
		"id":          doc.ID.Hex(),
		"chatId":      doc.ChatID,
		"senderId":    doc.SenderID,
		"receiverId":  doc.ReceiverID,
		"messageType": messageType,
		"text":        doc.Text,
		"imageUrl":    imageURL,
		"read":        doc.Read,
		"createdAt":   doc.CreatedAt,
	}
}

func validateChatParticipants(senderID, receiverID string) (string, int, string) {
	receiverID = strings.TrimSpace(receiverID)
	if receiverID == "" {
		return "", http.StatusBadRequest, "receiver_id is required"
	}
	if receiverID == senderID {
		return "", http.StatusBadRequest, "Cannot send message to yourself"
	}
	return receiverID, http.StatusOK, ""
}

func validateMessageLength(text string) (string, int, string) {
	text = strings.TrimSpace(text)
	if len([]rune(text)) > maxMessageLength {
		return "", http.StatusBadRequest, fmt.Sprintf("Message is too long. Maximum is %d characters", maxMessageLength)
	}
	return text, http.StatusOK, ""
}

func getChatUploadRoot() string {
	uploadRoot := strings.TrimSpace(os.Getenv("CHAT_UPLOAD_ROOT"))
	if uploadRoot == "" {
		uploadRoot = defaultChatUploadRoot
	}
	return uploadRoot
}

func generateRandomHex(length int) (string, error) {
	if length <= 0 {
		length = 8
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func detectImageExtension(contentType, fileName string) string {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "image/heic":
		return ".heic"
	case "image/heif":
		return ".heif"
	case "image/avif":
		return ".avif"
	}

	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(fileName)))
	switch ext {
	case ".jpg", ".jpeg":
		return ".jpg"
	case ".png", ".webp", ".gif", ".heic", ".heif", ".avif":
		return ext
	default:
		return ""
	}
}

func buildChatImagePublicURL(c *gin.Context, fileName string) string {
	publicBase := strings.TrimRight(strings.TrimSpace(os.Getenv("PUBLIC_BASE_URL")), "/")
	if publicBase != "" {
		return fmt.Sprintf("%s/uploads/%s/%s", publicBase, chatUploadSubDir, fileName)
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if forwardedProto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto")); forwardedProto != "" {
		scheme = strings.TrimSpace(strings.Split(forwardedProto, ",")[0])
	}

	return fmt.Sprintf("%s://%s/uploads/%s/%s", scheme, c.Request.Host, chatUploadSubDir, fileName)
}

func saveUploadedChatImage(c *gin.Context, fileHeader *multipart.FileHeader) (string, int, string) {
	if fileHeader == nil {
		return "", http.StatusBadRequest, "image file is required"
	}
	if fileHeader.Size <= 0 {
		return "", http.StatusBadRequest, "Image file is empty"
	}
	if fileHeader.Size > maxChatImageSizeBytes {
		maxMB := maxChatImageSizeBytes / (1024 * 1024)
		return "", http.StatusBadRequest, fmt.Sprintf("Image is too large. Maximum is %d MB", maxMB)
	}

	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", http.StatusBadRequest, "Failed to read uploaded image"
	}
	defer srcFile.Close()

	header := make([]byte, 512)
	bytesRead, readErr := srcFile.Read(header)
	if readErr != nil && readErr != io.EOF {
		return "", http.StatusBadRequest, "Failed to inspect uploaded image"
	}

	contentType := http.DetectContentType(header[:bytesRead])
	ext := detectImageExtension(contentType, fileHeader.Filename)
	if ext == "" {
		return "", http.StatusBadRequest, "Unsupported image format. Allowed: JPG, PNG, WEBP, GIF, HEIC"
	}

	uploadRoot := getChatUploadRoot()
	chatUploadsDir := filepath.Join(uploadRoot, chatUploadSubDir)
	if err := os.MkdirAll(chatUploadsDir, 0o755); err != nil {
		return "", http.StatusInternalServerError, "Failed to prepare upload directory"
	}

	randomPart, err := generateRandomHex(10)
	if err != nil {
		return "", http.StatusInternalServerError, "Failed to prepare image name"
	}
	fileName := fmt.Sprintf("%d_%s%s", time.Now().UTC().UnixNano(), randomPart, ext)
	fullPath := filepath.Join(chatUploadsDir, fileName)

	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		return "", http.StatusInternalServerError, "Failed to save uploaded image"
	}

	return buildChatImagePublicURL(c, fileName), http.StatusOK, ""
}

func GetChatMessages(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	requesterID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	otherUserID := parseChatPeerID(c)
	if otherUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if otherUserID == requesterID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot open chat with yourself"})
		return
	}

	limit := defaultChatLimit
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		n, convErr := strconv.Atoi(rawLimit)
		if convErr != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
			return
		}
		if n > maxChatLimit {
			n = maxChatLimit
		}
		limit = n
	}

	chatID := buildChatID(requesterID, otherUserID)

	filter := bson.M{
		"chat_id": chatID,
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := repositories.ChatMessagesCollection.Find(c.Request.Context(), filter, findOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat messages"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var docs []chatMessageDocument
	if err := cursor.All(c.Request.Context(), &docs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode chat messages"})
		return
	}

	result := make([]gin.H, 0, len(docs))
	for i := len(docs) - 1; i >= 0; i-- {
		result = append(result, toChatMessageResponse(docs[i]))
	}

	c.JSON(http.StatusOK, result)
}

func SendChatMessage(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	senderID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receiverID, status, validationError := validateChatParticipants(senderID, parseChatReceiverID(input))
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": validationError})
		return
	}

	text, status, validationError := validateMessageLength(asString(input["text"]))
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": validationError})
		return
	}
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message text cannot be empty"})
		return
	}

	now := time.Now().UTC()
	doc := chatMessageDocument{
		ChatID:      buildChatID(senderID, receiverID),
		SenderID:    senderID,
		ReceiverID:  receiverID,
		MessageType: messageTypeText,
		Text:        text,
		Read:        false,
		CreatedAt:   now,
	}

	insertResult, err := repositories.ChatMessagesCollection.InsertOne(c.Request.Context(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	insertedID, ok := insertResult.InsertedID.(bson.ObjectID)
	if ok {
		doc.ID = insertedID
	}

	// Send notification in background
	senderName := getUserDisplayNameByUID(c, senderID, nil)
	go func(senderName, receiverID, text string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		notifTitle := fmt.Sprintf("Новое сообщение от %s", senderName)
		notifLink := fmt.Sprintf("/chat?uid=%s", senderID)
		_ = CreateNotification(bgCtx, receiverID, "message", notifTitle, text, notifLink)
	}(senderName, receiverID, text)

	c.JSON(http.StatusCreated, toChatMessageResponse(doc))
}

func SendChatImageMessage(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	senderID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	receiverID, status, validationError := validateChatParticipants(senderID, parseChatReceiverIDFromForm(c))
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": validationError})
		return
	}

	text, status, validationError := validateMessageLength(c.PostForm("text"))
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": validationError})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	imageURL, status, imageError := saveUploadedChatImage(c, fileHeader)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": imageError})
		return
	}

	now := time.Now().UTC()
	doc := chatMessageDocument{
		ChatID:      buildChatID(senderID, receiverID),
		SenderID:    senderID,
		ReceiverID:  receiverID,
		MessageType: messageTypeImage,
		Text:        text,
		ImageURL:    imageURL,
		Read:        false,
		CreatedAt:   now,
	}

	insertResult, err := repositories.ChatMessagesCollection.InsertOne(c.Request.Context(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send image message"})
		return
	}

	insertedID, ok := insertResult.InsertedID.(bson.ObjectID)
	if ok {
		doc.ID = insertedID
	}

	// Send notification in background
	senderName := getUserDisplayNameByUID(c, senderID, nil)
	go func(senderName, receiverID, text string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		notifTitle := fmt.Sprintf("Новое сообщение от %s", senderName)
		notifLink := fmt.Sprintf("/chat?uid=%s", senderID)
		msgText := text
		if msgText == "" {
			msgText = "🖼️ Изображение"
		}
		_ = CreateNotification(bgCtx, receiverID, "message", notifTitle, msgText, notifLink)
	}(senderName, receiverID, text)

	c.JSON(http.StatusCreated, toChatMessageResponse(doc))
}

func MarkChatAsRead(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	requesterID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otherUserID := strings.TrimSpace(asString(input["user_id"]))
	if otherUserID == "" {
		otherUserID = strings.TrimSpace(asString(input["userId"]))
	}

	if otherUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	if otherUserID == requesterID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot mark own messages as read"})
		return
	}

	filter := bson.M{
		"chat_id":     buildChatID(requesterID, otherUserID),
		"sender_id":   otherUserID,
		"receiver_id": requesterID,
		"read":        false,
	}
	update := bson.M{
		"$set": bson.M{
			"read": true,
		},
	}

	result, err := repositories.ChatMessagesCollection.UpdateMany(c.Request.Context(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark chat messages as read"})
		return
	}

	// Also mark matching Firestore message notifications as read in background
	go func(receiverID, senderID string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		iter := repositories.FirestoreClient.Collection("notifications").
			Where("user_id", "==", receiverID).
			Where("type", "==", "message").
			Where("link", "==", fmt.Sprintf("/chat?uid=%s", senderID)).
			Where("read", "==", false).
			Documents(bgCtx)

		batch := repositories.FirestoreClient.Batch()
		count := 0
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				break
			}
			batch.Set(doc.Ref, map[string]interface{}{"read": true}, firestore.MergeAll)
			count++
			if count >= 500 {
				_, _ = batch.Commit(bgCtx)
				batch = repositories.FirestoreClient.Batch()
				count = 0
			}
		}
		if count > 0 {
			_, _ = batch.Commit(bgCtx)
		}
	}(requesterID, otherUserID)

	c.JSON(http.StatusOK, gin.H{
		"matched":  result.MatchedCount,
		"modified": result.ModifiedCount,
	})
}

func GetChatUnreadCount(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	requesterID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	filter := bson.M{
		"receiver_id": requesterID,
		"read":        false,
	}

	count, err := repositories.ChatMessagesCollection.CountDocuments(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate unread messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
