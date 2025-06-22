package middleware

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	verifier *oidc.IDTokenVerifier
)

func InitOIDC() {
	provider, err := oidc.NewProvider(context.Background(), os.Getenv("KEYCLOAK_ISSUER"))
	if err != nil {
		log.Fatalf("❌ Failed to create OIDC provider: %v", err)
	}

	verifier = provider.Verifier(&oidc.Config{
		ClientID: os.Getenv("KEYCLOAK_CLIENT_ID"),
	})
	log.Println("✅ OIDC initialized with Keycloak")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		idToken, err := verifier.Verify(c, rawToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims"})
			return
		}

		// 🔍 Keycloak stores roles inside "realm_access.roles"
		role := extractRole(claims)

		// 🚀 Set user info into context
		c.Set("user_id", claims["sub"])
		c.Set("username", claims["preferred_username"])
		c.Set("role", role)

		c.Next()
	}
}

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "ROLE_ADMIN" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
		c.Next()
	}
}

func extractRole(claims map[string]interface{}) string {
	realmAccess, ok := claims["realm_access"].(map[string]interface{})
	if !ok {
		return ""
	}

	roles, ok := realmAccess["roles"].([]interface{})
	if !ok {
		return ""
	}

	for _, r := range roles {
		if roleStr, ok := r.(string); ok && strings.HasPrefix(roleStr, "ROLE_") {
			return roleStr
		}
	}
	return ""
}
