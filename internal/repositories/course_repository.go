package repositories

import (
	"Internship/internal/entities"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *entities.Course) error
	GetByID(id uint) (*entities.Course, error)
	Update(course *entities.Course) error
	DeleteByID(id uint) error
}
type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *entities.Course) error {
	return r.db.Create(course).Error
}
func (r *courseRepository) GetByID(id uint) (*entities.Course, error) {
	var course entities.Course
	if err := r.db.First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}
func (r *courseRepository) Update(course *entities.Course) error {
	return r.db.Save(course).Error
}
func (r *courseRepository) DeleteByID(id uint) error {
	return r.db.Delete(&entities.Course{}, id).Error
}
