package repositories

import (
	"Internship/internal/entities"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Create(attachment *entities.Attachment) (*entities.Attachment, error)
	GetByID(id uint) (*entities.Attachment, error)
}

type attachmentRepo struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepo{db: db}
}

func (r *attachmentRepo) Create(attachment *entities.Attachment) (*entities.Attachment, error) {
	if err := r.db.Create(attachment).Error; err != nil {
		return nil, err
	}
	return attachment, nil
}

func (r *attachmentRepo) GetByID(id uint) (*entities.Attachment, error) {
	var attachment entities.Attachment
	if err := r.db.First(&attachment, id).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}
