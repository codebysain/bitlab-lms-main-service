package main

import (
	"Internship/internal/handler"
	"Internship/internal/middleware"
	"Internship/internal/repositories"
	"Internship/internal/service"
	"Internship/pkg/database"
	"Internship/pkg/minio"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

// @title BITLAB LMS API
// @version 1.0
// @description Backend service for BITLAB LMS platform
// @host localhost:8080
// @BasePath /
func main() {
	// Load .env for local/dev
	if err := godotenv.Load(); err != nil {
		logrus.Warn("⚠️ .env file not found (might be expected in Docker)")
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.Info("🚀 Starting Main Service...")

	// ✅ Init Keycloak OIDC verifier
	if os.Getenv("KEYCLOAK_ISSUER") == "" || os.Getenv("KEYCLOAK_CLIENT_ID") == "" {
		logrus.Fatal("❌ Missing Keycloak env (KEYCLOAK_ISSUER, KEYCLOAK_CLIENT_ID)")
	}
	middleware.InitOIDC()

	// Validate S3 env
	if os.Getenv("S3_BUCKET") == "" || os.Getenv("S3_ENDPOINT") == "" {
		logrus.Fatal("❌ Missing S3 envs (S3_BUCKET, S3_ENDPOINT)")
	}
	if os.Getenv("S3_ACCESS_KEY") == "" || os.Getenv("S3_SECRET_KEY") == "" {
		logrus.Fatal("❌ Missing S3 credentials (S3_ACCESS_KEY, S3_SECRET_KEY)")
	}
	logrus.Infof("📦 S3 bucket: %s", os.Getenv("S3_BUCKET"))

	// DB + MinIO
	db := database.Connect()
	logrus.Info("✅ Connected to database")
	minio.InitMinio()

	// Init Repos and Services
	userRepo := repositories.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)

	// Init Handlers
	courseHandler := handler.NewCourseHandler(service.NewCourseService(repositories.NewCourseRepository(db)))
	chapterHandler := handler.NewChapterHandler(service.NewChapterService(repositories.NewChapterRepository(db)))
	lessonHandler := handler.NewLessonHandler(service.NewLessonService(repositories.NewLessonRepository(db)))
	userHandler := handler.NewUserHandler(service.NewUserService(userRepo))
	authHandler := handler.NewAuthHandler(authService)
	refreshHandler := handler.NewRefreshHandler(authService)
	attachmentHandler := handler.NewAttachmentHandler(service.NewAttachmentService(repositories.NewAttachmentRepository(db)))

	// Init Router
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public
	router.POST("/login", authHandler.Login)
	router.POST("/refresh", refreshHandler.RefreshToken)

	// Protected routes
	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthMiddleware())

	// Admin-only
	authGroup.POST("/admin/register", middleware.AdminOnlyMiddleware(), authHandler.RegisterUser)
	authGroup.POST("/courses", middleware.AdminOnlyMiddleware(), courseHandler.CreateCourse)
	authGroup.POST("/chapters", middleware.AdminOnlyMiddleware(), chapterHandler.CreateChapter)
	authGroup.POST("/lessons", middleware.AdminOnlyMiddleware(), lessonHandler.CreateLesson)
	authGroup.POST("/upload", middleware.AdminOnlyMiddleware(), attachmentHandler.UploadFile)

	// Any authenticated user
	authGroup.PUT("/user/update", userHandler.UpdateUser)
	authGroup.GET("/courses/:id", courseHandler.GetCourseByID)
	authGroup.GET("/chapters/:id", chapterHandler.GetChapterByID)
	authGroup.GET("/lessons/:id", lessonHandler.GetLessonByID)
	authGroup.GET("/download/:id", attachmentHandler.DownloadFile)

	// Start
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("🌐 Server running on port %s", port)
	router.Run(":" + port)
}
