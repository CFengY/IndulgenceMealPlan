package repository

import (
	"IndulgenceMealPlan/model"
	"time"

	"gorm.io/gorm"
)

type MealRepository struct {
	db *gorm.DB
}

func NewMealRepository(db *gorm.DB) *MealRepository {
	return &MealRepository{db: db}
}

// Create 创建用餐记录
func (r *MealRepository) Create(meal *model.Meal) error {
	return r.db.Create(meal).Error
}

// GetByID 根据 MealID 获取用餐记录
func (r *MealRepository) GetByID(id uint) (*model.Meal, error) {
	var meal model.Meal
	if err := r.db.First(&meal, id).Error; err != nil {
		return nil, err
	}
	return &meal, nil
}

// GetMeals 根据用户 ID 获取所有用餐记录
func (r *MealRepository) GetMeals(userID uint) ([]model.Meal, error) {
	var meals []model.Meal
	err := r.db.Where("user_id = ?", userID).Order("meal_date DESC").Find(&meals).Error
	return meals, err
}

// Update 更新用餐记录
func (r *MealRepository) Update(meal *model.Meal) error {
	return r.db.Save(meal).Error
}

// Delete 删除用餐记录
func (r *MealRepository) Delete(id uint) error {
	return r.db.Delete(&model.Meal{}, id).Error
}

// ListByDateRange 获取指定日期范围内的用餐记录
func (r *MealRepository) ListByDateRange(userid uint, startDate, endDate time.Time) ([]model.Meal, error) {
	var meals []model.Meal
	err := r.db.Where("user_id = ? AND meal_date BETWEEN ? AND ?", userid, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("meal_date DESC").
		Find(&meals).Error
	return meals, err
}

// StatisticsByDateRange 获取指定日期范围内的食物统计数据
func (r *MealRepository) StatisticsByDateRange(userid uint, startDate, endDate time.Time) ([]model.FoodStatistic, error) {
	var stats []model.FoodStatistic
	err := r.db.Model(&model.Meal{}).
		Select("food_name, COUNT(*) as count").
		Where("user_id = ? AND meal_date BETWEEN ? AND ?", userid, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("food_name").
		Order("count DESC").
		Find(&stats).Error
	return stats, err
}
