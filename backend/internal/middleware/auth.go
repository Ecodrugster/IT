package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var authClient *auth.Client

func resolveServiceAccountPath() string {
	saPath := strings.TrimSpace(os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH"))
	if saPath != "" {
		if _, err := os.Stat(saPath); err == nil {
			return saPath
		}
		log.Printf("FIREBASE_SERVICE_ACCOUNT_PATH points to missing file: %s", saPath)
	}

	candidates := []string{
		"serviceAccountKey.json",
		"backend/serviceAccountKey.json",
		"../serviceAccountKey.json",
		"../../serviceAccountKey.json",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		// auth.go => internal/middleware, so go 2 levels up to backend root
		backendRoot := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "..", ".."))
		sourceCandidate := filepath.Join(backendRoot, "serviceAccountKey.json")
		if _, err := os.Stat(sourceCandidate); err == nil {
			return sourceCandidate
		}
	}

	return ""
}

func InitFirebase() {
	ctx := context.Background()
	saPath := resolveServiceAccountPath()
	firebaseJSON := strings.TrimSpace(os.Getenv("FIREBASE_JSON"))

	var app *firebase.App
	var err error

	if firebaseJSON != "" {
		opt := option.WithCredentialsJSON([]byte(firebaseJSON))
		app, err = firebase.NewApp(ctx, nil, opt)
		log.Printf("Firebase Auth initialized from FIREBASE_JSON environment variable")
	} else if saPath != "" {
		opt := option.WithCredentialsFile(saPath)
		app, err = firebase.NewApp(ctx, nil, opt)
		log.Printf("Firebase Auth initialized from credentials file: %s", saPath)
	} else {
		app, err = firebase.NewApp(ctx, nil)
		log.Printf("Firebase Auth initialized from default credentials")
	}

	if err != nil {
		log.Printf("error initializing firebase app: %v\n", err)
		return
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error getting firebase auth client: %v\n", err)
		return
	}

	authClient = client
	log.Println("Firebase Admin SDK initialized")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authClient == nil {
			c.Next() // Allow in dev if not configured
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		idToken := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store firebase UID in context
		c.Set("firebase_uid", token.UID)
		c.Next()
	}
}
