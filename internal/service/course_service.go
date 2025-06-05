package service

import (
	"Internship/internal/dto"
	"Internship/internal/entities"
	"Internship/internal/repositories"
	"context"
)

type CourseService interface {
	CreateCourse(ctx context.Context, req dto.CreateCourseRequest) (*entities.Course, error)
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

func (s *courseService) CreateCourse(ctx context.Context, req dto.CreateCourseRequest) (*entities.Course, error) {
	course := &entities.Course{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.repo.Create(ctx, course)
	return course, err
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
