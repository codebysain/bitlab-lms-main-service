package handler

import (
	"Internship/internal/entities"
	"Internship/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LessonHandler struct {
	service service.LessonService
}

func NewLessonHandler(service service.LessonService) *LessonHandler {
	return &LessonHandler{service: service}
}

// @Summary Get lesson by ID
// @Description Retrieves a lesson by its ID
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} entities.Lesson
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /lessons/{id} [get]
func (h *LessonHandler) GetLessonByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	lesson, err := h.service.GetLessonByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
		return
	}
	c.JSON(http.StatusOK, lesson)

}

// @Summary Create a new lesson
// @Description Creates a new lesson with name and description
// @Tags lessons
// @Accept json
// @Produce json
// @Param lesson body entities.Lesson true "Lesson to create"
// @Success 201 {object} entities.Lesson
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lessons [post]
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "ROLE_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	var lesson entities.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err := h.service.Create(&lesson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create lesson"})
		return
	}
	c.JSON(http.StatusCreated, lesson)
}
