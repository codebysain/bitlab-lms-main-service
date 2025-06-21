package main

import (
	"Internship/internal/handler"
	"Internship/internal/middleware"
	"Internship/internal/repositories"
	"Internship/internal/service"
	"Internship/pkg/database"
	"Internship/pkg/minio"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	// Validate required envs
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		logrus.Fatal("❌ Missing S3_BUCKET in env")
	}
	if os.Getenv("S3_ENDPOINT") == "" {
		logrus.Fatal("❌ Missing S3_ENDPOINT in env")
	}
	if os.Getenv("S3_ACCESS_KEY") == "" || os.Getenv("S3_SECRET_KEY") == "" {
		logrus.Fatal("❌ Missing MinIO credentials (S3_ACCESS_KEY or S3_SECRET_KEY)")
	}
	logrus.Infof("📦 Using S3 bucket: %s", bucket)

	// DB connect
	db := database.Connect()
	logrus.Info("✅ Connected to database")

	// Init MinIO
	minio.InitMinio()

	// Handlers
	courseHandler := handler.NewCourseHandler(service.NewCourseService(repositories.NewCourseRepository(db)))
	chapterHandler := handler.NewChapterHandler(service.NewChapterService(repositories.NewChapterRepository(db)))
	lessonHandler := handler.NewLessonHandler(service.NewLessonService(repositories.NewLessonRepository(db)))
	userRepo := repositories.NewUserRepository(db)
	userHandler := handler.NewUserHandler(service.NewUserService(userRepo))
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	refreshHandler := handler.NewRefreshHandler(authService)
	attachmentHandler := handler.NewAttachmentHandler(service.NewAttachmentService(repositories.NewAttachmentRepository(db)))

	// Router
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	router.POST("/login", authHandler.Login)
	router.POST("/refresh", refreshHandler.RefreshToken)

	// Protected routes
	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	authGroup.POST("/admin/register", authHandler.RegisterUser)
	authGroup.PUT("/user/update", userHandler.UpdateUser)
	authGroup.GET("/courses/:id", courseHandler.GetCourseByID)
	authGroup.POST("/courses", courseHandler.CreateCourse)
	authGroup.GET("/chapters/:id", chapterHandler.GetChapterByID)
	authGroup.POST("/chapters", chapterHandler.CreateChapter)
	authGroup.GET("/lessons/:id", lessonHandler.GetLessonByID)
	authGroup.POST("/lessons", lessonHandler.CreateLesson)
	authGroup.POST("/upload", attachmentHandler.UploadFile)
	authGroup.GET("/download/:id", attachmentHandler.DownloadFile)

	// Start
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("🌐 Server running on port %s", port)
	router.Run(":" + port)
}
