package repositories

import (
	"Internship/internal/entities"
	"errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*entities.User, error)
	Create(user *entities.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByUsername(username string) (*entities.User, error) {
	var user entities.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
func (r *userRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}
