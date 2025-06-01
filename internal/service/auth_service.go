package service

import (
	"Internship/internal/entities"
	"Internship/internal/repositories"
	"Internship/pkg/jwtutils"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Authenticate(username, password string) (string, string, error)
	GenerateTokens(user *entities.User) (string, string, error)
	RefreshTokens(refreshToken string) (string, string, error)
	RegisterUser(username, email, password, role string) error
	CreateUser(user *entities.User) error
}

type authService struct {
	userRepo repositories.UserRepository
	// Add jwtManager or key handling structs here if you split token logic
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Authenticate(username, password string) (string, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid password")
	}

	return s.GenerateTokens(user)
}

func (s *authService) GenerateTokens(user *entities.User) (string, string, error) {
	accessToken, err := jwtutils.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwtutils.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
func (s *authService) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := jwtutils.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return s.GenerateTokens(&entities.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	})
}

func (s *authService) RegisterUser(username, email, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entities.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	return s.userRepo.Create(user)
}
func (s *authService) CreateUser(user *entities.User) error {
	return s.userRepo.Create(user)
}
