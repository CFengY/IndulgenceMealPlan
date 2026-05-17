package handler

import (
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/service"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MealHandler struct {
	svc service.IMealService
}

func NewMealHandler(svc service.IMealService) *MealHandler {
	return &MealHandler{svc: svc}
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// Create 创建新用餐记录
func (h *MealHandler) Create(c *gin.Context) {
	mealType, err := strconv.Atoi(c.PostForm("meal_type"))
	if err != nil || mealType < 1 || mealType > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "meal_type 必须为 1-4（1早餐 2午餐 3晚餐 4零食）"})
		return
	}

	foodName := c.PostForm("food_name")
	if foodName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "food_name 不能为空"})
		return
	}

	mealDateStr := c.PostForm("meal_date")
	mealDate, err := time.Parse("2006-01-02", mealDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "meal_date 格式错误，应为 2006-01-02"})
		return
	}

	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	meal := &model.Meal{
		FoodName: foodName,
		MealType: model.MealType(mealType).String(),
		MealDate: mealDate,
		UserID:   currentId.(uint),
		Calories: parseFloat(c.PostForm("calories")),
		ProteinG: parseFloat(c.PostForm("protein_g")),
		FatG:     parseFloat(c.PostForm("fat_g")),
		CarbsG:   parseFloat(c.PostForm("carbs_g")),
	}

	if err := h.svc.Create(c, meal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": meal})
}

// GetMeals 根据用户id获取所有用餐记录
func (h *MealHandler) GetMeals(c *gin.Context) {
	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	meals, err := h.svc.GetMeals(c, currentId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": meals})
}

// Update 更新用餐记录
func (h *MealHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	meal, err := h.svc.GetByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	if meal.UserID != currentId.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改该记录"})
		return
	}

	if foodName := c.PostForm("food_name"); foodName != "" {
		meal.FoodName = foodName
	}
	if mealTypeStr := c.PostForm("meal_type"); mealTypeStr != "" {
		mealType, err := strconv.Atoi(mealTypeStr)
		if err != nil || mealType < 1 || mealType > 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "meal_type 必须为 1-4"})
			return
		}
		meal.MealType = model.MealType(mealType).String()
	}
	if mealDateStr := c.PostForm("meal_date"); mealDateStr != "" {
		mealDate, err := time.Parse("2006-01-02", mealDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "meal_date 格式错误"})
			return
		}
		meal.MealDate = mealDate
	}
	if caloriesStr := c.PostForm("calories"); caloriesStr != "" {
		meal.Calories = parseFloat(caloriesStr)
	}
	if proteinStr := c.PostForm("protein_g"); proteinStr != "" {
		meal.ProteinG = parseFloat(proteinStr)
	}
	if fatStr := c.PostForm("fat_g"); fatStr != "" {
		meal.FatG = parseFloat(fatStr)
	}
	if carbsStr := c.PostForm("carbs_g"); carbsStr != "" {
		meal.CarbsG = parseFloat(carbsStr)
	}

	if err := h.svc.Update(c, meal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": meal})
}

// Delete 删除用餐记录
func (h *MealHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	if err := h.svc.Delete(c, currentId.(uint), uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在或无权限删除"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ListByDateRange 根据日期范围获取用餐记录列表
func (h *MealHandler) ListByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date 和 end_date 不能为空"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date 格式错误，应为 2006-01-02"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date 格式错误，应为 2006-01-02"})
		return
	}

	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	userId := currentId.(uint)

	result, err := h.svc.ListByDateRange(c, userId, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ExportCSV 导出用餐记录为 CSV
func (h *MealHandler) ExportCSV(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date 和 end_date 不能为空"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date 格式错误"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date 格式错误"})
		return
	}

	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	result, err := h.svc.ListByDateRange(c, currentId.(uint), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("meals_export_%s_%s.csv", startDateStr, endDateStr)
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// BOM for Excel UTF-8 compatibility
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"日期", "食物名称", "餐类", "热量(kcal)", "蛋白质(g)", "脂肪(g)", "碳水(g)"})

	for _, meal := range result.Records {
		writer.Write([]string{
			meal.MealDate.Format("2006-01-02"),
			meal.FoodName,
			meal.MealType,
			fmt.Sprintf("%.1f", meal.Calories),
			fmt.Sprintf("%.1f", meal.ProteinG),
			fmt.Sprintf("%.1f", meal.FatG),
			fmt.Sprintf("%.1f", meal.CarbsG),
		})
	}

	writer.Flush()
}
