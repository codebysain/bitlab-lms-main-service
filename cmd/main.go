package main

import (
	"Internship/pkg/database"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Logging setup
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.Info("Starting Main Service...")

	// Load environment or set default
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=bitlab_lms port=5432 sslmode=disable"
	}

	// Connect to DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	logrus.Info("Connected to database")

	// Setup Gin
	r := gin.Default()

	// Ping route just to check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("Server running on port %s", port)
	r.Run(":" + port)
	db := database.Connect()

}
