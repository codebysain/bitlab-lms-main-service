package entities

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"default:'student'" json:"role"` // student, mentor, admin
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
