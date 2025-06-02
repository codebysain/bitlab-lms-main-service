package service

import (
	"context"
	"errors"

	"Internship/internal/dto"
	"Internship/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	UpdateUser(ctx context.Context, userID uint, req dto.UpdateUserRequest) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) UpdateUser(ctx context.Context, userID uint, req dto.UpdateUserRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Username = req.Username
	user.Email = req.Email
	user.Password = string(hashedPassword)

	return s.repo.Update(ctx, user)
}
