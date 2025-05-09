package service_test

import (
	"Internship/internal/entities"
	"Internship/internal/mocks"
	"Internship/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCourseByID_Success(t *testing.T) {
	mockRepo := new(mocks.CourseRepository)
	courseSvc := service.NewCourseService(mockRepo)

	expected := &entities.Course{
		ID:          1,
		Name:        "Test Course",
		Description: "A mock test course",
	}
	mockRepo.On("GetByID", uint(1)).Return(expected, nil)

	result, err := courseSvc.GetCourseByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetCourseByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.CourseRepository)
	courseSvc := service.NewCourseService(mockRepo)

	mockRepo.On("GetByID", uint(99)).Return(nil, assert.AnError)

	result, err := courseSvc.GetCourseByID(99)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateCourse_Success(t *testing.T) {
	mockRepo := new(mocks.CourseRepository)
	courseSvc := service.NewCourseService(mockRepo)

	newCourse := &entities.Course{
		Name:        "Go for Juniors",
		Description: "Starter course for new Go devs",
	}

	mockRepo.On("Create", newCourse).Return(nil)

	err := courseSvc.Create(newCourse)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateCourse_Failure(t *testing.T) {
	mockRepo := new(mocks.CourseRepository)
	courseSvc := service.NewCourseService(mockRepo)

	newCourse := &entities.Course{
		Name:        "Fail Course",
		Description: "This will trigger an error",
	}

	mockRepo.On("Create", newCourse).Return(assert.AnError)

	err := courseSvc.Create(newCourse)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCourse(t *testing.T) {
	mockRepo := new(mocks.CourseRepository)
	courseSvc := service.NewCourseService(mockRepo)

	course := &entities.Course{
		ID:          1,
		Name:        "Updated Course",
		Description: "Updated desc",
	}

	mockRepo.On("Update", course).Return(nil)

	err := courseSvc.Update(course)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCourseByID(t *testing.T) {
	mockRepo := new(mocks.CourseRepository)
	courseSvc := service.NewCourseService(mockRepo)

	mockRepo.On("DeleteByID", uint(1)).Return(nil)

	err := courseSvc.DeleteCourseByID(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
