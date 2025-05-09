package handler_test

import (
	"Internship/internal/entities"
	"Internship/internal/handler"
	"Internship/internal/mocks"
	"Internship/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCourseByIDHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.CourseRepository)
	svc := service.NewCourseService(mockRepo)
	h := handler.NewCourseHandler(svc)

	course := &entities.Course{
		ID:          1,
		Name:        "Unit Test Course",
		Description: "For handler testing",
	}
	mockRepo.On("GetByID", uint(1)).Return(course, nil)

	// Setup router and request
	router := gin.Default()
	router.GET("/courses/:id", h.GetCourseByID)

	req := httptest.NewRequest(http.MethodGet, "/courses/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	mockRepo.AssertExpectations(t)
}
func TestGetCourseByIDHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.CourseRepository)
	svc := service.NewCourseService(mockRepo)
	h := handler.NewCourseHandler(svc)

	mockRepo.On("GetByID", uint(99)).Return(nil, assert.AnError)

	router := gin.Default()
	router.GET("/courses/:id", h.GetCourseByID)

	req := httptest.NewRequest(http.MethodGet, "/courses/99", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	mockRepo.AssertExpectations(t)
}
