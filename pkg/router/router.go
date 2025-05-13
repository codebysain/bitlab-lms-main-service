package router

import (
	"Internship/internal/handler"
	"Internship/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine,
	courseHandler *handler.CourseHandler,
	chapterHandler *handler.ChapterHandler,
	lessonHandler *handler.LessonHandler,
) {
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
}
