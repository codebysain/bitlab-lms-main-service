package router

import (
	"Internship/internal/handler"
	"Internship/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/*
Public →   /healthz, /swagger/*, /login, /refresh
Protected → (AuthMiddleware)
  - Any role      : GET course / chapter / lesson, PUT user, file download
  - Admin only    : register user, create course/chapter/lesson, file upload
*/
func SetupRoutes(
	r *gin.Engine,
	authH *handler.AuthHandler,
	refreshH *handler.RefreshHandler,
	userH *handler.UserHandler,
	courseH *handler.CourseHandler,
	chapterH *handler.ChapterHandler,
	lessonH *handler.LessonHandler,
	attachH *handler.AttachmentHandler,
) {
	/* ---------- public ---------- */
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/login", authH.Login)
	r.POST("/refresh", refreshH.RefreshToken)

	/* ---------- protected ---------- */
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// admin-only subgroup
	admin := protected.Group("/")
	admin.Use(middleware.AdminOnlyMiddleware())
	{
		admin.POST("/admin/register", authH.RegisterUser)
		admin.POST("/courses", courseH.CreateCourse)
		admin.POST("/chapters", chapterH.CreateChapter)
		admin.POST("/lessons", lessonH.CreateLesson)
		admin.POST("/upload", attachH.UploadFile)
	}

	// all roles
	protected.PUT("/user/update", userH.UpdateUser)
	protected.GET("/courses/:id", courseH.GetCourseByID)
	protected.GET("/chapters/:id", chapterH.GetChapterByID)
	protected.GET("/lessons/:id", lessonH.GetLessonByID)
	protected.GET("/download/:id", attachH.DownloadFile)

	/* ---------- 404 fallback ---------- */
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "route not found"})
	})
}
