package service

import (
	"IndulgenceMealPlan/config"
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/repository"
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	global.Config = &config.Config{
		Jwt: config.JwtConfig{
			Secret:     "test-secret",
			Expiration: 3600,
			Name:       "Authorization",
		},
		Upload: config.UploadConfig{
			Dir:     "./testdata/uploads",
			MaxSize: 10 << 20,
		},
	}

	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer mr.Close()
	global.Redis = redis.NewClient(&redis.Options{Addr: mr.Addr()})

	gin.SetMode(gin.TestMode)

	os.MkdirAll("./testdata/uploads", 0755)
	code := m.Run()
	os.RemoveAll("./testdata")
	os.Exit(code)
}

func newTestGinContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", nil)
	return c
}

func TestMealService_GetByID_NoCache(t *testing.T) {
	mockRepo := new(repository.MockMealRepository)
	expectedMeal := &model.Meal{FoodName: "测试食物"}
	expectedMeal.ID = 1
	mockRepo.On("GetByID", uint(1)).Return(expectedMeal, nil)

	svc := NewMealService(mockRepo, "./testdata", 10<<20)

	// 清除 Redis 缓存
	global.Redis.Del(context.Background(), "meal:1")

	meal, err := svc.GetByID(newTestGinContext(), 1)
	require.NoError(t, err)
	assert.Equal(t, "测试食物", meal.FoodName)
	mockRepo.AssertExpectations(t)
}

func TestMealService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(repository.MockMealRepository)
	mockRepo.On("GetByID", uint(999)).Return(nil, assert.AnError)

	svc := NewMealService(mockRepo, "./testdata", 10<<20)

	global.Redis.Del(context.Background(), "meal:999")

	_, err := svc.GetByID(newTestGinContext(), 999)
	assert.Error(t, err)
}

func TestMealService_GetMeals_Success(t *testing.T) {
	mockRepo := new(repository.MockMealRepository)
	expectedMeals := []model.Meal{
		{FoodName: "早餐", MealType: "早餐"},
		{FoodName: "午餐", MealType: "午餐"},
	}
	mockRepo.On("GetMeals", uint(1)).Return(expectedMeals, nil)

	svc := NewMealService(mockRepo, "./testdata", 10<<20)

	global.Redis.Del(context.Background(), "meals:user:1")

	meals, err := svc.GetMeals(newTestGinContext(), 1)
	require.NoError(t, err)
	assert.Len(t, meals, 2)
	mockRepo.AssertExpectations(t)
}

func TestMealService_GetMealsNoCache(t *testing.T) {
	mockRepo := new(repository.MockMealRepository)
	expectedMeals := []model.Meal{{FoodName: "直接查询"}}
	mockRepo.On("GetMeals", uint(1)).Return(expectedMeals, nil)

	svc := NewMealService(mockRepo, "./testdata", 10<<20)

	meals, err := svc.GetMealsNoCache(1)
	require.NoError(t, err)
	assert.Len(t, meals, 1)
}

func TestMealService_GetByIDNoCache(t *testing.T) {
	mockRepo := new(repository.MockMealRepository)
	expectedMeal := &model.Meal{FoodName: "无缓存查询"}
	expectedMeal.ID = 1
	mockRepo.On("GetByID", uint(1)).Return(expectedMeal, nil)

	svc := NewMealService(mockRepo, "./testdata", 10<<20)

	meal, err := svc.GetByIDNoCache(1)
	require.NoError(t, err)
	assert.Equal(t, "无缓存查询", meal.FoodName)
}

func TestMealService_ListByDateRange_Success(t *testing.T) {
	mockRepo := new(repository.MockMealRepository)
	records := []model.Meal{{FoodName: "范围查询食物"}}
	stats := []model.FoodStatistic{{FoodName: "范围查询食物", Count: 1}}

	mockRepo.On("ListByDateRange", uint(1), mock.Anything, mock.Anything).Return(records, nil)
	mockRepo.On("StatisticsByDateRange", uint(1), mock.Anything, mock.Anything).Return(stats, nil)

	svc := NewMealService(mockRepo, "./testdata", 10<<20)

	startDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)

	global.Redis.Del(context.Background(), "meals:user:1:2026-05-01:2026-05-12")

	result, err := svc.ListByDateRange(newTestGinContext(), 1, startDate, endDate)
	require.NoError(t, err)
	assert.Len(t, result.Records, 1)
	assert.Len(t, result.Statistics, 1)
	mockRepo.AssertExpectations(t)
}

func TestMealService_Delete_PermissionDenied(t *testing.T) {
	mockRepo := new(repository.MockMealRepository)
	otherUserMeal := &model.Meal{FoodName: "别人的食物", UserID: 2}
	otherUserMeal.ID = 1
	mockRepo.On("GetByID", uint(1)).Return(otherUserMeal, nil)

	svc := NewMealService(mockRepo, "./testdata", 10<<20)

	global.Redis.Del(context.Background(), "meal:1")

	err := svc.Delete(newTestGinContext(), 1, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无权限")
}

func TestMealService_Delete_Success(t *testing.T) {
	mockRepo := new(repository.MockMealRepository)
	meal := &model.Meal{FoodName: "我的食物", UserID: 1}
	meal.ID = 1
	mockRepo.On("GetByID", uint(1)).Return(meal, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	svc := NewMealService(mockRepo, "./testdata", 10<<20)

	global.Redis.Del(context.Background(), "meal:1")

	err := svc.Delete(newTestGinContext(), 1, 1)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
