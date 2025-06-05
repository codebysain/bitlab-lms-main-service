package middleware

import (
	"Internship/pkg/jwtutils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		const BearerSchema = "Bearer "

		if !strings.HasPrefix(authHeader, BearerSchema) {
			log.Println("❌ Missing or malformed Authorization header:", authHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, BearerSchema)
		if len(tokenStr) > 50 {
			log.Println("📦 Extracted Token:", tokenStr[:50]+"...")
		} else {
			log.Println("📦 Extracted Token:", tokenStr)
		}

		claims, err := jwtutils.ParseAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		log.Printf("🎯 Extracted Claims: %+v\n", claims)

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
