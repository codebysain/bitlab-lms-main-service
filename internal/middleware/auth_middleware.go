package middleware

import (
	"Internship/pkg/jwtutils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		log.Println("📦 Extracted Token:", tokenStr[:50]+"...") // don't print full token, just first 50 chars

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			log.Println("🔍 Token Header:", token.Header)

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				log.Println("❌ Unexpected signing method:", token.Header["alg"])
				return nil, errors.New("unexpected signing method")
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				log.Println("❌ No kid in token header")
				return nil, errors.New("kid header not found")
			}

			log.Println("🔑 kid:", kid)

			pubKey, err := jwtutils.FetchPublicKey(kid)
			if err != nil {
				log.Println("❌ Failed to fetch public key:", err)
			}
			return pubKey, err
		})

		if err != nil {
			log.Println("❌ Token parse error:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			log.Println("❌ Token is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("❌ Invalid token claims format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		log.Println("🧾 Token Claims:", claims)

		realmAccess, ok := claims["realm_access"].(map[string]interface{})
		if !ok {
			log.Println("⚠️ No realm_access in claims")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No realm access info in token"})
			return
		}

		rawRoles, ok := realmAccess["roles"].([]interface{})
		if !ok {
			log.Println("⚠️ Roles not found in realm_access")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Roles missing in token"})
			return
		}

		log.Println("✅ Roles extracted:", rawRoles)

		c.Set("roles", rawRoles)
		c.Next()
	}
}
