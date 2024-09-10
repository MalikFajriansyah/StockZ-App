package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string    `json:"username" gorm:"size:255;unique;not null"`
	Email     string    `json:"email" gorm:"size:255;not null; unique"`
	Password  string    `json:"password" gorm:"not null"`
	Posts     []Post    `gorm:"foreignKey:UserID"`
	Comments  []Comment `gorm:"foreignKey:UserID"`
	Likes     []Like    `gorm:"foreignKey:UserID"`
	Followers []Follow  `gorm:"foreignKey:FollowedID"`
	Following []Follow  `gorm:"foreignKey:FollowerID"`
}
