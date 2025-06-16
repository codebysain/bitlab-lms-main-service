package handler

import (
	"Internship/internal/service"
	"Internship/pkg/minio"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AttachmentHandler struct {
	service service.AttachmentService
}

func NewAttachmentHandler(s service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: s}
}

func (h *AttachmentHandler) UploadFile(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || (role != "ROLE_ADMIN" && role != "ROLE_TEACHER") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	lessonIDStr := c.PostForm("lesson_id")
	lessonID, err := strconv.Atoi(lessonIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File open failed"})
		return
	}
	defer src.Close()

	uploadedURL, hashedName, err := minio.UploadFile(file.Filename, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload to MinIO failed"})
		return
	}

	attachment, err := h.service.CreateAttachment(file.Filename, hashedName, uploadedURL, uint(lessonID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB save failed"})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func (h *AttachmentHandler) DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
		return
	}

	attachment, err := h.service.GetAttachmentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	fileReader, err := minio.DownloadFile(attachment.NameHashed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Download failed"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+attachment.NameOriginal)
	c.DataFromReader(http.StatusOK, -1, "application/octet-stream", fileReader, nil)
}
