package repositories

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client

func InitFirestore() {
	ctx := context.Background()
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	
	if projectID == "" {
		projectID = "itstep-social" // Fallback
	}

	// Initialize Firebase App
	// In production, use GOOGLE_APPLICATION_CREDENTIALS env var
	var err error
	var opt option.ClientOption

	// 1. Проверяем наличие JSON в переменной окружения (для Railway)
	firebaseJSON := os.Getenv("FIREBASE_JSON")
	if firebaseJSON != "" {
		opt = option.WithCredentialsJSON([]byte(firebaseJSON))
		log.Println("Firebase initialized from FIREBASE_JSON environment variable")
	} else {
		// 2. Локальный запуск - используем путь из .env или дефолтный файл
		serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
		if serviceAccountPath == "" {
			serviceAccountPath = "serviceAccountKey.json"
		}
		opt = option.WithCredentialsFile(serviceAccountPath)
		log.Printf("Firebase initialized from file: %s", serviceAccountPath)
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	FirestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore client: %v\n", err)
	}
	
	log.Println("Firestore connection established")
}

var (
	Descending = firestore.Desc
	MergeAll   = firestore.MergeAll
)
