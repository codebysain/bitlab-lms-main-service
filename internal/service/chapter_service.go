package service

import (
	"Internship/internal/entities"
	"Internship/internal/repositories"
)

type ChapterService interface {
	Create(chapter *entities.Chapter) error
	GetChapterByID(chapterID uint) (*entities.Chapter, error)
	Update(chapter *entities.Chapter) error
	DeleteChapterByID(chapterID uint) error
}
type chapterService struct {
	repo repositories.ChapterRepository
}

func NewChapterService(repo repositories.ChapterRepository) ChapterService {
	return &chapterService{repo: repo}
}

func (s *chapterService) Create(chapter *entities.Chapter) error {
	return s.repo.Create(chapter)
}
func (s *chapterService) GetChapterByID(chapterID uint) (*entities.Chapter, error) {
	return s.repo.GetByID(chapterID)
}
func (s *chapterService) Update(chapter *entities.Chapter) error {
	return s.repo.Update(chapter)
}
func (s *chapterService) DeleteChapterByID(chapterID uint) error {
	return s.repo.DeleteByID(chapterID)
}
