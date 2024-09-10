package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content" gorm:"size:500;not null"`
	UserID  uint   `json:"user_id" gorm:"not null"` // Foreign key to User
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PostID  uint   `json:"post_id" gorm:"not null"` // Foreign key to Post
	Post    Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
