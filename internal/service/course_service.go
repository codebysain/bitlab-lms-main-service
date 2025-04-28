package service

import (
	"Internship/internal/entities"
	"Internship/internal/repositories"
)

type CourseService interface {
	Create(course *entities.Course) error
	GetCourseByID(courseID uint) (*entities.Course, error)
	Update(course *entities.Course) error
	DeleteCourseByID(courseID uint) error
}

type courseService struct {
	repo repositories.CourseRepository
}

func NewCourseService(repo repositories.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) Create(course *entities.Course) error {
	return s.repo.Create(course)
}
func (s *courseService) GetCourseByID(courseID uint) (*entities.Course, error) {
	return s.repo.GetByID(courseID)
}
func (s *courseService) Update(course *entities.Course) error {
	return s.repo.Update(course)
}
func (s *courseService) DeleteCourseByID(courseID uint) error {
	return s.repo.DeleteByID(courseID)
}
