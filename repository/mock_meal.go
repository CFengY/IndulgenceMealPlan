package repository

import (
	"IndulgenceMealPlan/model"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockMealRepository struct {
	mock.Mock
}

func (m *MockMealRepository) Create(meal *model.Meal) error {
	args := m.Called(meal)
	return args.Error(0)
}

func (m *MockMealRepository) GetByID(id uint) (*model.Meal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Meal), args.Error(1)
}

func (m *MockMealRepository) GetMeals(userID uint) ([]model.Meal, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Meal), args.Error(1)
}

func (m *MockMealRepository) Update(meal *model.Meal) error {
	args := m.Called(meal)
	return args.Error(0)
}

func (m *MockMealRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMealRepository) ListByDateRange(userid uint, startDate, endDate time.Time) ([]model.Meal, error) {
	args := m.Called(userid, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Meal), args.Error(1)
}

func (m *MockMealRepository) StatisticsByDateRange(userid uint, startDate, endDate time.Time) ([]model.FoodStatistic, error) {
	args := m.Called(userid, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.FoodStatistic), args.Error(1)
}

func (m *MockMealRepository) NutritionByDateRange(userid uint, startDate, endDate time.Time) (*model.NutritionSummary, error) {
	args := m.Called(userid, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NutritionSummary), args.Error(1)
}
