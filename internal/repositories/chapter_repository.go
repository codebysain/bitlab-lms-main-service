package repositories

import (
	"Internship/internal/entities"
	"gorm.io/gorm"
)

type ChapterRepository interface {
	Create(chapter *entities.Chapter) error
	GetByID(chapterID uint) (*entities.Chapter, error)
	Update(chapter *entities.Chapter) error
	DeleteByID(chapterID uint) error
}

type chapterRepository struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) ChapterRepository {
	return &chapterRepository{db: db}
}
func (r *chapterRepository) Create(chapter *entities.Chapter) error {
	return r.db.Create(chapter).Error
}
func (r *chapterRepository) GetByID(chapterID uint) (*entities.Chapter, error) {
	var chapter entities.Chapter
	if err := r.db.Preload("Lessons").First(&chapter, chapterID).Error; err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (r *chapterRepository) Update(chapter *entities.Chapter) error {
	return r.db.Save(chapter).Error
}
func (r *chapterRepository) DeleteByID(chapterID uint) error {
	return r.db.Delete(&entities.Chapter{}, chapterID).Error
}
