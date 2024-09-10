package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"not null"` // Foreign key to User
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PostID uint `json:"post_id" gorm:"not null"` // Foreign key to Post
	Post   Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
