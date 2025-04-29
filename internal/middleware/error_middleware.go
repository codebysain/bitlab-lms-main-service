package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			statusCode := http.StatusInternalServerError
			message := "Internal Server Error"
			for _, err := range c.Errors {
				if err.Type == gin.ErrorTypePublic {
					statusCode = http.StatusNotFound
					message = err.Error()
					break
				}
			}
			c.JSON(statusCode, gin.H{"error": message})
			c.Abort()
		}
	}
}
