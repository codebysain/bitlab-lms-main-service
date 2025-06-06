package handler

import (
	"Internship/internal/entities"
	"Internship/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ChapterHandler struct {
	service service.ChapterService
}

func NewChapterHandler(service service.ChapterService) *ChapterHandler {
	return &ChapterHandler{service: service}
}

// GetChapterByID godoc
// @Summary      Get chapter by ID
// @Description  Get a chapter by its ID
// @Tags         chapters
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Chapter ID"
// @Success      200  {object}  entities.Chapter
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /chapters/{id} [get]
func (h *ChapterHandler) GetChapterByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	chapter, err := h.service.GetChapterByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	c.JSON(http.StatusOK, chapter)
}

// CreateChapter godoc
// @Summary      Create a new chapter
// @Description  Creates a new chapter
// @Tags         chapters
// @Accept       json
// @Produce      json
// @Param        chapter body entities.Chapter true "Chapter to create"
// @Success      201  {object}  entities.Chapter
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /chapters [post]
func (h *ChapterHandler) CreateChapter(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "ROLE_ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	var chapter entities.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.Create(&chapter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chapter"})
		return
	}

	c.JSON(http.StatusCreated, chapter)
}
