package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	Content string `gorm:"type:text;not null" json:"content"`
	Images  string `gorm:"type:text" json:"images,omitempty"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (*Post) TableName() string {
	return "posts"
}
