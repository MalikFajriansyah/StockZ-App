package model

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	FollowerID uint `json:"follower_id" gorm:"not null"` // Foreign key to User (the follower)
	Follower   User `gorm:"foreignKey:FollowerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FollowedID uint `json:"followed_id" gorm:"not null"` // Foreign key to User (the followed)
	Followed   User `gorm:"foreignKey:FollowedID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
