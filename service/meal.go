package service

import (
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/repository"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type IMealService interface {
	Create(c *gin.Context, meal *model.Meal) error
	GetByID(c *gin.Context, id uint) (*model.Meal, error)
	GetByIDNoCache(id uint) (*model.Meal, error)
	GetMeals(c *gin.Context, userID uint) ([]model.Meal, error)
	GetMealsNoCache(userID uint) ([]model.Meal, error)
	Update(c *gin.Context, meal *model.Meal) error
	Delete(c *gin.Context, userid, id uint) error
	ListByDateRange(c *gin.Context, userid uint, startDate, endDate time.Time) (*DateRangeResult, error)
}

type MealService struct {
	repo      *repository.MealRepository
	uploadDir string
	maxSize   int64
}

func NewMealService(repo *repository.MealRepository, uploadDir string, maxSize int64) IMealService {
	return &MealService{repo: repo, uploadDir: uploadDir, maxSize: maxSize}
}

// Create 创建新用餐记录，支持图片上传
func (s *MealService) Create(c *gin.Context, meal *model.Meal) error {
	file, err := c.FormFile("image")
	if err == nil {
		if file.Size > s.maxSize {
			return fmt.Errorf("图片大小超过限制（最大 %d 字节）", s.maxSize)
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			return fmt.Errorf("不支持的图片格式: %s", ext)
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join(s.uploadDir, filename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			return fmt.Errorf("保存图片失败: %w", err)
		}

		meal.ImagePath = filename
	}

	ctx := c.Request.Context()
	script := redis.NewScript(`
		-- 删除主键
		redis.call("DEL", "meals:user:" .. KEYS[1])

		-- 扫描并删除所有相关键
		local cursor = "0"
		repeat
			local res = redis.call("SCAN", cursor, "MATCH", "meals:user:" .. KEYS[1] .. ":*", "COUNT", 100)
			cursor = res[1]
			local keys = res[2]
			
			if #keys > 0 then
				redis.call("DEL", unpack(keys))
			end
		until cursor == "0"

		return true
	`)

	err = script.Run(
		ctx,
		global.Redis,
		[]string{
			strconv.FormatUint(uint64(meal.UserID), 10),
		},
	).Err()

	if err != nil {
		log.Println("缓存删除失败：", err)
	}

	return s.repo.Create(meal)
}

// GetByID 根据 MealID 获取用餐记录
func (s *MealService) GetByID(c *gin.Context, id uint) (*model.Meal, error) {
	ctx := c.Request.Context()

	if str := global.Redis.Get(ctx, "meal:"+fmt.Sprint(id)).Val(); str != "" {
		err := json.Unmarshal([]byte(str), &model.Meal{})
		if err == nil {
			return &model.Meal{}, nil
		}
	}

	meal, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(meal)
	if err != nil {
		return nil, err
	}

	global.Redis.Set(ctx, "meal:"+fmt.Sprint(id), jsonData, time.Minute*10)

	return meal, nil
}

// GetByID 根据 MealID 获取用餐记录 （不使用缓存，直接从数据库查询）
func (s *MealService) GetByIDNoCache(id uint) (*model.Meal, error) {

	return s.repo.GetByID(id)
}

// GetMeals 根据用户 ID 获取所有用餐记录
func (s *MealService) GetMeals(c *gin.Context, userID uint) ([]model.Meal, error) {
	ctx := c.Request.Context()

	if str := global.Redis.Get(ctx, "meals:user:"+fmt.Sprint(userID)).Val(); str != "" {
		var meals []model.Meal
		err := json.Unmarshal([]byte(str), &meals)
		if err == nil {
			return meals, nil
		}
	}

	meals, err := s.repo.GetMeals(userID)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(meals)
	if err != nil {
		return nil, err
	}

	global.Redis.Set(ctx, "meals:user:"+fmt.Sprint(userID), jsonData, time.Minute*10)

	return meals, nil
}

// GetMeals 根据用户 ID 获取所有用餐记录 （不使用缓存，直接从数据库查询）
func (s *MealService) GetMealsNoCache(userID uint) ([]model.Meal, error) {
	return s.repo.GetMeals(userID)
}

// Update 更新用餐记录，支持图片替换
func (s *MealService) Update(c *gin.Context, meal *model.Meal) error {
	file, err := c.FormFile("image")
	if err == nil {
		if file.Size > s.maxSize {
			return fmt.Errorf("图片大小超过限制（最大 %d 字节）", s.maxSize)
		}

		// 删除旧图片
		if meal.ImagePath != "" {
			oldPath := filepath.Join(s.uploadDir, meal.ImagePath)
			os.Remove(oldPath)
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join(s.uploadDir, filename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			return fmt.Errorf("保存图片失败: %w", err)
		}

		meal.ImagePath = filename
	}

	ctx := c.Request.Context()
	script := redis.NewScript(`
		redis.call("DEL", "meal:" .. KEYS[1])
		redis.call("DEL", "meals:user:" .. KEYS[2])

		local cursor = "0"
		repeat
			local res = redis.call("SCAN", cursor, "MATCH", "meals:user:" .. KEYS[2] .. ":*", "COUNT", 100)
			cursor = res[1]
			for _, key in ipairs(res[2]) do
				redis.call("DEL", key)
			end
		until cursor == "0"

		return true
	`)

	err = script.Run(
		ctx,
		global.Redis,
		[]string{
			strconv.FormatUint(uint64(meal.ID), 10),
			strconv.FormatUint(uint64(meal.UserID), 10),
		},
	).Err()

	if err != nil {
		log.Println("缓存删除失败：", err)
	}

	return s.repo.Update(meal)
}

// Delete 删除用餐记录，同时删除关联图片
func (s *MealService) Delete(c *gin.Context, userid, id uint) error {

	meal, err := s.GetByID(c, id)

	if err != nil {
		return err
	}

	if meal.UserID != userid {
		return fmt.Errorf("无权限删除该记录")
	}

	ctx := c.Request.Context()

	script := redis.NewScript(`
		redis.call("DEL", "meal:" .. KEYS[1])
		redis.call("DEL", "meals:user:" .. KEYS[2])

		local cursor = "0"
		repeat
			local res = redis.call("SCAN", cursor, "MATCH", "meals:user:" .. KEYS[2] .. ":*", "COUNT", 100)
			cursor = res[1]
			for _, key in ipairs(res[2]) do
				redis.call("DEL", key)
			end
		until cursor == "0"

		return true
	`)

	err = script.Run(
		ctx,
		global.Redis,
		[]string{
			strconv.FormatUint(uint64(id), 10),
			strconv.FormatUint(uint64(userid), 10),
		},
	).Err()

	if err != nil {
		log.Println("缓存删除失败：", err)
	}

	// global.Redis.Del(ctx, "meal:"+fmt.Sprint(id))
	// global.Redis.Del(ctx, "meals:user:"+fmt.Sprint(userid))

	// 删除关联图片

	if meal.ImagePath != "" {
		os.Remove(filepath.Join(s.uploadDir, meal.ImagePath))
	}

	return s.repo.Delete(id)
}

type DateRangeResult struct {
	Records    []model.Meal          `json:"records"`
	Statistics []model.FoodStatistic `json:"statistics"`
}

// ListByDateRange 获取指定日期范围内的用餐记录和食物统计数据
func (s *MealService) ListByDateRange(c *gin.Context, userid uint, startDate, endDate time.Time) (*DateRangeResult, error) {
	ctx := c.Request.Context()

	if str := global.Redis.Get(ctx, string(fmt.Sprintf("meals:user:%d:%s:%s", userid, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))).Val(); str != "" {
		var result DateRangeResult
		err := json.Unmarshal([]byte(str), &result)
		if err == nil {
			return &result, nil
		}
	}

	records, err := s.repo.ListByDateRange(userid, startDate, endDate)
	if err != nil {
		return nil, err
	}

	stats, err := s.repo.StatisticsByDateRange(userid, startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := &DateRangeResult{Records: records, Statistics: stats}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	global.Redis.Set(ctx, string(fmt.Sprintf("meals:user:%d:%s:%s", userid, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))), jsonData, time.Minute*10)

	return result, nil
}
