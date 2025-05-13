package main

import (
	"Internship/internal/handler"
	"Internship/internal/repositories"
	"Internship/internal/service"
	"Internship/pkg/database"
	"Internship/pkg/router"
	"github.com/joho/godotenv"
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

	// 🔥 Register routes
	router.SetupRoutes(r, courseHandler, chapterHandler, lessonHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("Server running on port %s", port)
	r.Run(":" + port)
}
