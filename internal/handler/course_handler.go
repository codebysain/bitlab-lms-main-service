package handler

import (
	"Internship/internal/dto"
	"Internship/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(service service.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

// GetCourseByID godoc
// @Summary      Get course by ID
// @Description  Get a course by its ID
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  entities.Course
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	course, err := h.service.GetCourseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// @Summary Create a new course
// @Description Creates a new course with name and description
// @Tags courses
// @Accept json
// @Produce json
// @Param course body entities.Course true "Course to create"
// @Success 201 {object} entities.Course
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /courses [post]
// CreateCourse creates a new course
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "ROLE_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var req dto.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.service.CreateCourse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
	log.Println("🧠 ROLE in context:", role)

}
