package main

import (
	"Internship/internal/handler"
	"Internship/internal/middleware"
	"Internship/internal/repositories"
	"Internship/internal/service"
	"Internship/pkg/database"
	"Internship/pkg/minio"
	"Internship/pkg/router" // <- new import
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

/*
@title       BITLAB LMS API
@version     1.0
@description Backend service for BITLAB LMS platform
@host        localhost:8080
@BasePath    /
*/
func main() {
	_ = godotenv.Load() // ignore error inside Docker
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	requireEnv(
		"KEYCLOAK_ISSUER", "KEYCLOAK_CLIENT_ID", "KEYCLOAK_TOKEN_URL",
		"S3_BUCKET", "S3_ENDPOINT", "S3_ACCESS_KEY", "S3_SECRET_KEY",
	)

	logrus.Info("🚀 Starting BITLAB-LMS main service")

	/* ---------- infrastructure ---------- */
	middleware.InitOIDC()
	db := database.Connect()
	minio.InitMinio()

	/* ---------- repos / services ---------- */
	userRepo := repositories.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)

	courseH := handler.NewCourseHandler(service.NewCourseService(repositories.NewCourseRepository(db)))
	chapterH := handler.NewChapterHandler(service.NewChapterService(repositories.NewChapterRepository(db)))
	lessonH := handler.NewLessonHandler(service.NewLessonService(repositories.NewLessonRepository(db)))
	attachH := handler.NewAttachmentHandler(service.NewAttachmentService(repositories.NewAttachmentRepository(db)))

	authH := handler.NewAuthHandler(authSvc)
	refreshH := handler.NewRefreshHandler(authSvc)
	userH := handler.NewUserHandler(service.NewUserService(userRepo))

	/* ---------- router ---------- */
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.ErrorHandlerMiddleware())
	router.SetupRoutes(r, authH, refreshH, userH, courseH, chapterH, lessonH, attachH)

	/* ---------- run w/ graceful shutdown ---------- */
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logrus.Infof("🌐 Listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("❌ listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("🛑 Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("❌ forced shutdown: %v", err)
	}
	logrus.Info("✅ Server exited gracefully")
}

/* helpers */
func requireEnv(keys ...string) {
	for _, k := range keys {
		if os.Getenv(k) == "" {
			logrus.Fatalf("❌ required env %s not set", k)
		}
	}
}
