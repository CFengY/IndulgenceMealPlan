package handler

import (
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMealService struct {
	mock.Mock
}

func (m *MockMealService) Create(c *gin.Context, meal *model.Meal) error {
	args := m.Called(c, meal)
	return args.Error(0)
}

func (m *MockMealService) GetByID(c *gin.Context, id uint) (*model.Meal, error) {
	args := m.Called(c, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Meal), args.Error(1)
}

func (m *MockMealService) GetByIDNoCache(id uint) (*model.Meal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Meal), args.Error(1)
}

func (m *MockMealService) GetMeals(c *gin.Context, userID uint) ([]model.Meal, error) {
	args := m.Called(c, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Meal), args.Error(1)
}

func (m *MockMealService) GetMealsNoCache(userID uint) ([]model.Meal, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Meal), args.Error(1)
}

func (m *MockMealService) Update(c *gin.Context, meal *model.Meal) error {
	args := m.Called(c, meal)
	return args.Error(0)
}

func (m *MockMealService) Delete(c *gin.Context, userid, id uint) error {
	args := m.Called(c, userid, id)
	return args.Error(0)
}

func (m *MockMealService) ListByDateRange(c *gin.Context, userid uint, startDate, endDate time.Time) (*service.DateRangeResult, error) {
	args := m.Called(c, userid, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.DateRangeResult), args.Error(1)
}

func setupTestContext(method, path string, body url.Values) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, path, strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	return c, w
}

func TestMealHandler_Create_Success(t *testing.T) {
	mockSvc := new(MockMealService)
	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("*model.Meal")).Return(nil)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/meals", url.Values{
		"food_name": {"测试食物"},
		"meal_type": {"1"},
		"meal_date": {"2026-05-12"},
	})
	c.Set("currentId", uint(1))

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestMealHandler_Create_MissingFoodName(t *testing.T) {
	mockSvc := new(MockMealService)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/meals", url.Values{
		"meal_type": {"1"},
		"meal_date": {"2026-05-12"},
	})

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMealHandler_Create_InvalidMealType(t *testing.T) {
	mockSvc := new(MockMealService)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/meals", url.Values{
		"food_name": {"测试"},
		"meal_type": {"5"},
		"meal_date": {"2026-05-12"},
	})

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMealHandler_Create_NoCurrentId(t *testing.T) {
	mockSvc := new(MockMealService)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/meals", url.Values{
		"food_name": {"测试"},
		"meal_type": {"1"},
		"meal_date": {"2026-05-12"},
	})

	handler.Create(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMealHandler_GetMeals_Success(t *testing.T) {
	mockSvc := new(MockMealService)
	expectedMeals := []model.Meal{
		{FoodName: "早餐A", MealType: "早餐"},
	}
	mockSvc.On("GetMeals", mock.Anything, uint(1)).Return(expectedMeals, nil)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodGet, "/api/v1/meals", url.Values{})
	c.Set("currentId", uint(1))

	handler.GetMeals(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestMealHandler_Delete_Success(t *testing.T) {
	mockSvc := new(MockMealService)
	mockSvc.On("Delete", mock.Anything, uint(1), uint(1)).Return(nil)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodDelete, "/api/v1/meals/1", url.Values{})
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("currentId", uint(1))

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestMealHandler_Delete_InvalidID(t *testing.T) {
	mockSvc := new(MockMealService)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodDelete, "/api/v1/meals/abc", url.Values{})
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMealHandler_Update_NotFound(t *testing.T) {
	mockSvc := new(MockMealService)
	mockSvc.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("record not found"))
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodPut, "/api/v1/meals/999", url.Values{})
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	handler.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMealHandler_Update_Forbidden(t *testing.T) {
	existingMeal := &model.Meal{UserID: 2, FoodName: "别人的食物"}
	mockSvc := new(MockMealService)
	mockSvc.On("GetByID", mock.Anything, uint(1)).Return(existingMeal, nil)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodPut, "/api/v1/meals/1", url.Values{})
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("currentId", uint(1))

	handler.Update(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMealHandler_ListByDateRange_Success(t *testing.T) {
	mockSvc := new(MockMealService)
	result := &service.DateRangeResult{
		Records:    []model.Meal{},
		Statistics: []model.FoodStatistic{},
	}
	mockSvc.On("ListByDateRange", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(result, nil)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodGet, "/api/v1/meals/range?start_date=2026-05-01&end_date=2026-05-12", url.Values{})
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/meals/range?start_date=2026-05-01&end_date=2026-05-12", nil)
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Set("currentId", uint(1))

	handler.ListByDateRange(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestMealHandler_ListByDateRange_MissingDates(t *testing.T) {
	mockSvc := new(MockMealService)
	handler := NewMealHandler(mockSvc)

	c, w := setupTestContext(http.MethodGet, "/api/v1/meals/range", url.Values{})
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/meals/range", nil)

	handler.ListByDateRange(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
