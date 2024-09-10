package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string    `json:"tittle" gorm:"size:255; not null"`
	Description string    `json:"description" gorm:"size:500; not null"`
	MediaURL    string    `json:"media_url" gorm:"size:500;not null"`
	MediaType   string    `json:"media_type" gorm:"size:50;not null"` // Example: "image/jpeg", "video/mp4"
	UserID      uint      `json:"user_id" gorm:"not null"`            // Foreign key to User
	User        User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Comments    []Comment `gorm:"foreignKey:PostID"`
	Likes       []Like    `gorm:"foreignKey:PostID"`
}
