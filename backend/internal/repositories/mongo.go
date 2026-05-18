package repositories

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	MongoClient            *mongo.Client
	ChatMessagesCollection *mongo.Collection
)

func InitMongo() {
	mongoURL := strings.TrimSpace(os.Getenv("MONGO_URL"))
	if mongoURL == "" {
		log.Println("MONGO_URL is not set. Chat API is disabled.")
		return
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Printf("Failed to initialize Mongo client: %v", err)
		return
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		_ = client.Disconnect(context.Background())
		return
	}

	dbName := strings.TrimSpace(os.Getenv("MONGO_DB_NAME"))
	if dbName == "" {
		dbName = "itstep_chat"
	}

	MongoClient = client
	ChatMessagesCollection = client.Database(dbName).Collection("messages")
	log.Printf("MongoDB connection established (db=%s, collection=messages)", dbName)

	indexCtx, indexCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer indexCancel()

	_, err = ChatMessagesCollection.Indexes().CreateMany(indexCtx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "chat_id", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "receiver_id", Value: 1},
				{Key: "read", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to create Mongo indexes for chat messages: %v", err)
	}
}

func IsMongoChatReady() bool {
	return ChatMessagesCollection != nil
}
