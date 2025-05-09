package main

import (
	"Internship/internal/handler"
	"Internship/internal/middleware"
	"Internship/internal/repositories"
	"Internship/internal/service"
	"Internship/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Internship/docs" // for Swagger
)

// @title BITLAB LMS API
// @version 1.0
// @description Backend service for BITLAB LMS platform
// @host localhost:8080
// @BasePath /
func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.Info("Starting Main Service...")

	db := database.Connect()
	logrus.Info("Connected to database")

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

	r := gin.Default()
	r.Use(middleware.ErrorHandlerMiddleware())

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	courseGroup := r.Group("/courses")
	{
		courseGroup.GET("/:id", courseHandler.GetCourseByID)
		courseGroup.POST("", courseHandler.CreateCourse)
	}
	chapterGroup := r.Group("/chapters")
	{
		chapterGroup.GET("/:id", chapterHandler.GetChapterByID)
		chapterGroup.POST("", chapterHandler.CreateChapter)
	}
	lessonGroup := r.Group("/lessons")
	{
		lessonGroup.GET("/:id", lessonHandler.GetLessonByID)
		lessonGroup.POST("", lessonHandler.CreateLesson)
	}

	logrus.Info("Server running on port 8080")
	r.Run(":8080")
}
