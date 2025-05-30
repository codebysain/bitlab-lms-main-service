package service

import (
	"Internship/internal/entities"
	"Internship/internal/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Authenticate(username, password string) (string, string, error)
	GenerateTokens(user *entities.User) (string, string, error)
	RefreshTokens(refreshToken string) (string, string, error)
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

// Dummy token gen – replace with real JWT generation
func (s *authService) GenerateTokens(user *entities.User) (string, string, error) {
	accessToken := "fake-access-token-for-" + user.Username
	refreshToken := "fake-refresh-token-for-" + user.Username
	return accessToken, refreshToken, nil
}
func (s *authService) RefreshTokens(refreshToken string) (string, string, error) {
	// Dummy logic for now: always return new fake tokens
	accessToken := "new-fake-access-token"
	newRefreshToken := "new-fake-refresh-token"

	return accessToken, newRefreshToken, nil
}
