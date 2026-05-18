package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/user/itstep-backend/internal/handlers"
	"github.com/user/itstep-backend/internal/middleware"
	"github.com/user/itstep-backend/internal/repositories"
)

func loadEnv() {
	candidates := []string{
		".env",
		"../.env",
		"../../.env",
		"backend/.env",
	}

	for _, candidate := range candidates {
		if err := godotenv.Load(candidate); err == nil {
			log.Printf("Loaded environment from %s", candidate)
			return
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("No .env file found, using defaults/environment variables")
		cwd = ""
	}

	dir := cwd
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(candidate); statErr == nil {
			if loadErr := godotenv.Load(candidate); loadErr == nil {
				log.Printf("Loaded environment from %s", candidate)
				return
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fallback for `go run` binaries executed from temp dirs:
	// use the source file location embedded in debug symbols.
	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		mainDir := filepath.Dir(sourceFile)
		sourceCandidates := []string{
			filepath.Join(mainDir, "..", "..", ".env"),
			filepath.Join(mainDir, "..", ".env"),
		}

		for _, candidate := range sourceCandidates {
			cleanCandidate := filepath.Clean(candidate)
			if _, statErr := os.Stat(cleanCandidate); statErr == nil {
				if loadErr := godotenv.Load(cleanCandidate); loadErr == nil {
					log.Printf("Loaded environment from %s", cleanCandidate)
					return
				}
			}
		}
	}

	log.Printf("No .env file found (cwd: %s), using defaults/environment variables", cwd)
}

func main() {
	// Load .env file from common run locations
	loadEnv()

	// Initialize Firestore
	repositories.InitFirestore()
	// Initialize MongoDB for chats
	repositories.InitMongo()

	// Initialize Firebase
	middleware.InitFirebase()

	r := gin.Default()

	// CORS Setup
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	uploadRoot := strings.TrimSpace(os.Getenv("CHAT_UPLOAD_ROOT"))
	if uploadRoot == "" {
		uploadRoot = "uploads"
	}
	if err := os.MkdirAll(filepath.Join(uploadRoot, "chat"), 0o755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
	}
	r.Static("/uploads", filepath.Clean(uploadRoot))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		v1.GET("/profile", handlers.GetUserProfile)
		v1.PUT("/profile", handlers.UpdateUserProfile)
		v1.POST("/bootstrap/admin", handlers.BootstrapFirstAdmin)

		// Posts
		v1.POST("/posts", handlers.CreatePost)
		v1.GET("/posts", handlers.GetPosts)
		v1.POST("/posts/:id/like", handlers.LikePost)
		v1.POST("/posts/:id/comments", handlers.AddComment)
		v1.GET("/posts/:id/comments", handlers.GetComments)

		// News
		v1.GET("/news", handlers.GetNews)

		// User stats
		v1.GET("/profile/stats", handlers.GetUserStats)
		// Users
		v1.GET("/users", handlers.GetAllUsers)

		// Clubs
		v1.GET("/clubs", handlers.GetClubs)
		v1.POST("/clubs", handlers.CreateClub)
		v1.PUT("/clubs/:id", handlers.UpdateClub)
		v1.DELETE("/clubs/:id", handlers.DeleteClub)
		v1.POST("/clubs/:id/join", handlers.JoinClub)
		v1.POST("/clubs/:id/leave", handlers.LeaveClub)
		v1.GET("/grades", handlers.GetUserGrades)
		v1.GET("/schedule", handlers.GetSchedule)

		// Chat (MongoDB)
		v1.GET("/chat/messages", handlers.GetChatMessages)
		v1.POST("/chat/messages", handlers.SendChatMessage)
		v1.POST("/chat/messages/image", handlers.SendChatImageMessage)
		v1.POST("/chat/messages/read", handlers.MarkChatAsRead)
		v1.GET("/chat/unread-count", handlers.GetChatUnreadCount)

		// Notifications
		v1.GET("/notifications", handlers.GetNotifications)
		v1.POST("/notifications/:id/read", handlers.MarkNotificationAsRead)
		v1.POST("/notifications/read-all", handlers.MarkAllNotificationsAsRead)
		v1.GET("/notifications/unread-count", handlers.GetUnifiedUnreadCount)

		// College Store
		v1.GET("/store/items", handlers.GetStoreItems)
		v1.POST("/store/purchase", handlers.PurchaseItem)
		v1.GET("/store/history", handlers.GetPurchaseHistory)

		// Teacher routes
		teacher := v1.Group("/teacher")
		teacher.Use(middleware.TeacherRequired())
		{
			teacher.POST("/grades", handlers.AddGrade)
			teacher.GET("/schedule", handlers.GetTeacherSchedule)
			teacher.POST("/attendance", handlers.MarkAttendance)
			teacher.GET("/attendance", handlers.GetAttendance)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("/store/purchases", handlers.GetAdminPurchases)
			admin.POST("/store/claim", handlers.ClaimPurchase)

			admin.GET("/dashboard", handlers.AdminGetDashboard)
			admin.GET("/users", handlers.AdminGetUsers)
			admin.PUT("/users/:id/role", handlers.AdminUpdateUserRole)
			admin.PUT("/users/:id/group", handlers.AdminUpdateUserGroup)
			admin.GET("/posts", handlers.AdminGetPosts)
			admin.DELETE("/posts/:id", handlers.AdminDeletePost)
			admin.POST("/news", handlers.CreateNews)
			admin.DELETE("/news/:id", handlers.AdminDeleteNews)
			admin.PUT("/news/:id", handlers.AdminUpdateNews)
			admin.GET("/clubs", handlers.AdminGetClubs)
			admin.POST("/clubs", handlers.AdminCreateClub)
			admin.DELETE("/clubs/:id", handlers.AdminDeleteClub)
			admin.PUT("/clubs/:id", handlers.AdminUpdateClub)
			admin.GET("/club-requests", handlers.AdminGetClubRequests)
			admin.POST("/club-requests/:id/approve", handlers.AdminApproveClubRequest)
			admin.POST("/club-requests/:id/reject", handlers.AdminRejectClubRequest)
			admin.GET("/schedule", handlers.AdminGetSchedule)
			admin.POST("/schedule", handlers.AdminCreateSchedule)
			admin.PUT("/schedule/:id", handlers.AdminUpdateSchedule)
			admin.DELETE("/schedule/:id", handlers.AdminDeleteSchedule)
		}
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
