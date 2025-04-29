package main

import (
	"Internship/internal/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.Info("Starting Main Service...")

	//	db := database.Connect()

	logrus.Info("Connected to database")

	r := gin.Default()
	r.Use(middleware.ErrorHandlerMiddleware())
	// Initialize repositories
	//	courseRepo := repositories.NewCourseRepository(db)
	//	chapterRepo := repositories.NewChapterRepository(db)
	//	lessonRepo := repositories.NewLessonRepository(db)

	// Create a Course
	//	course := &entities.Course{
	//		Name:        "Backend Development with Go",
	//		Description: "Learn how to build scalable backend systems with Golang",
	//	}
	//	if err := courseRepo.Create(course); err != nil {
	//		panic("failed to create course: " + err.Error())
	//	}
	//	fmt.Println("Course created with ID:", course.ID)

	// Create a Chapter
	//	chapter := &entities.Chapter{
	//		Name:        "GORM Basics",
	//		Description: "Learn ORM in Go",
	//		Order:       1,
	//		CourseID:    course.ID,
	//	}
	//	if err := chapterRepo.Create(chapter); err != nil {
	//		panic("failed to create chapter: " + err.Error())
	//	}
	//	fmt.Println("Chapter created with ID:", chapter.ID)
	//
	// Create a Lesson
	//	lesson := &entities.Lesson{
	//		Name:        "Connecting to DB with GORM",
	//		Description: "Setup and first queries",
	//		Content:     "Install GORM, connect to PostgreSQL...",
	//		Order:       1,
	//		ChapterID:   chapter.ID,
	//	}
	//	if err := lessonRepo.Create(lesson); err != nil {
	//		panic("failed to create lesson: " + err.Error())
	//	}
	//	fmt.Println("Lesson created with ID:", lesson.ID)

	r.GET("/test-error", func(c *gin.Context) {
		c.AbortWithError(404, fmt.Errorf("test course not found")).SetType(gin.ErrorTypePublic)
	})
	r.Run(":8080")

}
