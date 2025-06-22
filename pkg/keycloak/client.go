package keycloak

import (
	"context"
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

var (
	Client       = gocloak.NewClient(getKeycloakURL())
	Realm        = getEnv("KEYCLOAK_REALM", "bitlab")
	ClientID     = getEnv("KEYCLOAK_CLIENT_ID", "bitlab-lms")
	ClientSecret = getEnv("KEYCLOAK_CLIENT_SECRET", "")
)

func getKeycloakURL() string {
	url := os.Getenv("KEYCLOAK_ISSUER")
	if url == "" {
		url = "http://localhost:8080/realms/bitlab"
	}
	return url[:len(url)-len("/realms/"+Realm)] // strip /realms/bitlab → base URL
}

func getEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return fallback
}

func GetAdminToken(ctx context.Context) (*gocloak.JWT, error) {
	username := os.Getenv("KEYCLOAK_ADMIN")
	password := os.Getenv("KEYCLOAK_ADMIN_PASSWORD")

	token, err := Client.LoginClient(ctx, ClientID, ClientSecret, Realm)
	if err != nil {
		fmt.Println("❌ Client login failed, trying with username and password...")

		token, err = Client.LoginAdmin(ctx, username, password, Realm)
		if err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
	}
	return token, nil
}
