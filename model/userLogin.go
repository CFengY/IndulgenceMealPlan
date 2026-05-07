package model

import "gorm.io/gorm"

type User struct {
	// ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	gorm.Model
	Username string `gorm:"type:varchar(50);unique;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Meals    []Meal `json:"meals,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (*User) TableName() string {
	return "users"
}
