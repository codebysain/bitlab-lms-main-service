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

	"os"

	_ "Internship/docs" // for Swagger
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title BITLAB LMS API
// @version 1.0
// @description Backend service for BITLAB LMS platform
// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.Info("Starting Main Service...")

	db := database.Connect()
	logrus.Info("Connected to database")
	minio.InitMinio()

	// Init repos, services, handlers
	courseRepo := repositories.NewCourseRepository(db)
	courseService := service.NewCourseService(courseRepo)
	courseHandler := handler.NewCourseHandler(courseService)

	chapterRepo := repositories.NewChapterRepository(db)
	chapterService := service.NewChapterService(chapterRepo)
	chapterHandler := handler.NewChapterHandler(chapterService)

	lessonRepo := repositories.NewLessonRepository(db)
	lessonService := service.NewLessonService(lessonRepo)
	lessonHandler := handler.NewLessonHandler(lessonService)

	userRepo := repositories.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	refreshHandler := handler.NewRefreshHandler(authService)
	attachmentRepo := repositories.NewAttachmentRepository(db)
	attachmentService := service.NewAttachmentService(attachmentRepo)
	attachmentHandler := handler.NewAttachmentHandler(attachmentService)

	// Single router instance
	router := gin.Default()

	// Swagger
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("Server running on port %s", port)
	router.Run(":" + port)

}
