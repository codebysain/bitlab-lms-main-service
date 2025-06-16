package service

import (
	"Internship/internal/entities"
	"Internship/internal/repositories"
)

type AttachmentService interface {
	CreateAttachment(originalName, hashedName, url string, lessonID uint) (*entities.Attachment, error)
	GetAttachmentByID(id uint) (*entities.Attachment, error)
}

type attachmentService struct {
	repo repositories.AttachmentRepository
}

func NewAttachmentService(r repositories.AttachmentRepository) AttachmentService {
	return &attachmentService{repo: r}
}

func (s *attachmentService) CreateAttachment(originalName, hashedName, url string, lessonID uint) (*entities.Attachment, error) {
	attachment := &entities.Attachment{
		NameOriginal: originalName,
		NameHashed:   hashedName,
		URL:          url,
		LessonID:     lessonID,
	}
	return s.repo.Create(attachment)
}

func (s *attachmentService) GetAttachmentByID(id uint) (*entities.Attachment, error) {
	return s.repo.GetByID(id)
}
