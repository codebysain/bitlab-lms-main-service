package entities

import "time"

type Attachment struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NameOriginal string    `gorm:"not null" json:"name_original"`
	NameHashed   string    `gorm:"not null;unique" json:"name_hashed"`
	URL          string    `gorm:"not null" json:"url"`
	LessonID     uint      `json:"lesson_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
