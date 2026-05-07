package model

import (
	"time"

	"gorm.io/gorm"
)

type MealType int

const (
	MealTypeBreakfast MealType = 1
	MealTypeLunch     MealType = 2
	MealTypeDinner    MealType = 3
	// MealTypeSnack     MealType = 4
)

func (m MealType) String() string {
	switch m {
	case MealTypeBreakfast:
		return "早餐"
	case MealTypeLunch:
		return "午餐"
	case MealTypeDinner:
		return "晚餐"
	// case MealTypeSnack:
	// 	return "零食"
	default:
		return "未知"
	}
}

type Meal struct {

	// ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	gorm.Model
	UserID    uint      `json:"user_id"`
	FoodName  string    `gorm:"type:varchar(100);not null" json:"food_name"`
	MealType  string    `gorm:"type:varchar(10);not null" json:"meal_type"`
	MealDate  time.Time `gorm:"type:date;not null" json:"meal_date"`
	ImagePath string    `gorm:"type:varchar(255)" json:"image_path,omitempty"`
	// CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	// UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	// DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (*Meal) TableName() string {
	return "meals"
}

type FoodStatistic struct {
	FoodName string `json:"food_name"`
	Count    int64  `json:"count"`
}
