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
		if authHeader == "" {
			log.Println("❌ Missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("❌ Invalid Authorization header format: %s\n", authHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("🔐 Verifying token: %s\n", rawToken)

		idToken, err := verifier.Verify(c, rawToken)
		if err != nil {
			log.Printf("❌ Token verification failed: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token verification failed", "details": err.Error()})
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			log.Printf("❌ Failed to parse token claims: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims", "details": err.Error()})
			return
		}

		log.Printf("✅ Token verified. Claims: %+v\n", claims)

		role := extractRole(claims)
		if role == "" {
			log.Println("⚠️ No ROLE_ found in realm_access.roles")
		} else {
			log.Printf("🛡️ Extracted role: %s\n", role)
		}

		c.Set("user_id", claims["sub"])
		c.Set("username", claims["preferred_username"])
		c.Set("role", role)

		c.Next()
	}
}

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			log.Println("❌ No role found in context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role missing in context"})
			return
		}

		if role != "ROLE_ADMIN" {
			log.Printf("⛔ Access denied. Role found: %v\n", role)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		log.Println("✅ Admin role confirmed. Proceeding.")
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
