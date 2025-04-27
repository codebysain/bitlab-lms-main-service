package repositories

import (
	"Internship/internal/entities"
	"gorm.io/gorm"
)

type LessonRepository interface {
	Create(lesson *entities.Lesson) error
	GetByID(id uint) (*entities.Lesson, error)
	Update(lesson *entities.Lesson) error
	DeleteByID(id uint) error
}
type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db: db}
}
func (r *lessonRepository) Create(lesson *entities.Lesson) error {
	return r.db.Create(lesson).Error
}
func (r *lessonRepository) GetByID(id uint) (*entities.Lesson, error) {
	var lesson entities.Lesson
	if err := r.db.First(&lesson, id).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}
func (r *lessonRepository) Update(lesson *entities.Lesson) error {
	return r.db.Save(lesson).Error
}
func (r *lessonRepository) DeleteByID(id uint) error {
	return r.db.Delete(&entities.Lesson{}, id).Error
}
