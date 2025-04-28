package service

import (
	"Internship/internal/entities"
	"Internship/internal/repositories"
)

type LessonService interface {
	Create(lesson *entities.Lesson) error
	GetLessonByID(lessonID uint) (*entities.Lesson, error)
	Update(lesson *entities.Lesson) error
	DeleteLessonByID(lessonID uint) error
}

type lessonService struct {
	repo repositories.LessonRepository
}

func NewLessonService(repo repositories.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func (s *lessonService) Create(lesson *entities.Lesson) error {
	return s.repo.Create(lesson)
}
func (s *lessonService) GetLessonByID(lessonID uint) (*entities.Lesson, error) {
	return s.repo.GetByID(lessonID)
}
func (s *lessonService) Update(lesson *entities.Lesson) error {
	return s.repo.Update(lesson)
}
func (s *lessonService) DeleteLessonByID(lessonID uint) error {
	return s.repo.DeleteByID(lessonID)
}
