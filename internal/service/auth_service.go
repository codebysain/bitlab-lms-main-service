package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"Internship/internal/entities"
	"Internship/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Authenticate(username, password string) (access, refresh string, err error)
	RefreshTokens(ctx context.Context, refresh string) (access, newRefresh string, err error)
	RegisterUser(username, email, password, role string) error
	CreateUser(user *entities.User) error
}

type authService struct {
	userRepo     repositories.UserRepository
	tokenURL     string
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenURL:     mustEnv("KEYCLOAK_TOKEN_URL"),
		clientID:     mustEnv("KEYCLOAK_CLIENT_ID"),
		clientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"), // empty if public client
		httpClient:   &http.Client{Timeout: 5 * time.Second},
	}
}

/*
---------------------------------------------------

	LOGIN (username / password)

----------------------------------------------------
*/
func (s *authService) Authenticate(username, password string) (string, string, error) {
	// Optional local DB check (leave it if you still store users)
	user, err := s.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		return "", "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid password")
	}

	form := url.Values{
		"grant_type": {"password"},
		"client_id":  {s.clientID},
		"username":   {username},
		"password":   {password},
	}
	if s.clientSecret != "" {
		form.Set("client_secret", s.clientSecret)
	}

	tok, err := s.postToKeycloak(context.Background(), form)
	if err != nil {
		return "", "", err
	}
	return tok.AccessToken, tok.RefreshToken, nil
}

/*
---------------------------------------------------

	REFRESH FLOW

----------------------------------------------------
*/
func (s *authService) RefreshTokens(ctx context.Context, refresh string) (string, string, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {s.clientID},
		"refresh_token": {refresh},
	}
	if s.clientSecret != "" {
		form.Set("client_secret", s.clientSecret)
	}

	tok, err := s.postToKeycloak(ctx, form)
	if err != nil {
		return "", "", err
	}
	return tok.AccessToken, tok.RefreshToken, nil
}

/*
---------------------------------------------------

	USER REG / DB HELPERS (unchanged)

----------------------------------------------------
*/
func (s *authService) RegisterUser(username, email, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &entities.User{
		Username: username,
		Email:    email,
		Password: string(hash),
		Role:     role,
	}
	return s.userRepo.Create(user)
}

func (s *authService) CreateUser(user *entities.User) error {
	return s.userRepo.Create(user)
}

/*
---------------------------------------------------

	INTERNAL HELPERS

----------------------------------------------------
*/
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *authService) postToKeycloak(ctx context.Context, form url.Values) (*tokenResponse, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, s.tokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak error %d: %s", resp.StatusCode, body)
	}

	var t tokenResponse
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env %s not set", key))
	}
	return v
}
