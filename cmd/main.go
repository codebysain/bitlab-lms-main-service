package main

import (
	"Internship/internal/entities"
	"Internship/internal/middleware"
	"Internship/internal/repositories"
	"Internship/pkg/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.Info("Starting Main Service...")

	db := database.Connect()

	logrus.Info("Connected to database")

	r := gin.Default()
	r.Use(middleware.ErrorHandlerMiddleware())
	logrus.Info("Creating new course...")
	courseRepo := repositories.NewCourseRepository(db)
	chapterRepo := repositories.NewChapterRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)

	course := &entities.Course{
		Name:        "Backend Development with Go",
		Description: "Learn how to build scalable backend systems with Golang",
	}

	logrus.Debugf("Course data: %+v", course)

	if err := courseRepo.Create(course); err != nil {
		logrus.Errorf("Failed to create course: %v", err)
		panic("failed to create course: " + err.Error())
	}

	logrus.Infof("Course created with ID: %d", course.ID)

	logrus.Info("Creating new chapter...")

	chapter := &entities.Chapter{
		Name:        "GORM Basics",
		Description: "Learn ORM in Go",
		Order:       1,
		CourseID:    course.ID,
	}

	logrus.Debugf("Chapter data: %+v", chapter)

	if err := chapterRepo.Create(chapter); err != nil {
		logrus.Errorf("Failed to create chapter: %v", err)
		panic("failed to create chapter: " + err.Error())
	}

	logrus.Infof("Chapter created with ID: %d", chapter.ID)

	logrus.Info("Creating new lesson...")

	lesson := &entities.Lesson{
		Name:        "Connecting to DB with GORM",
		Description: "Setup and first queries",
		Content:     "Install GORM, connect to PostgreSQL...",
		Order:       1,
		ChapterID:   chapter.ID,
	}

	logrus.Debugf("Lesson data: %+v", lesson)

	if err := lessonRepo.Create(lesson); err != nil {
		logrus.Errorf("Failed to create lesson: %v", err)
		panic("failed to create lesson: " + err.Error())
	}

	logrus.Infof("Lesson created with ID: %d", lesson.ID)

	r.GET("/test-error", func(c *gin.Context) {
		c.AbortWithError(404, fmt.Errorf("test course not found")).SetType(gin.ErrorTypePublic)
	})
	r.Run(":8080")

}
