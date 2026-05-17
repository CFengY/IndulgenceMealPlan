package handler

import (
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc service.IPostService
}

func NewPostHandler(svc service.IPostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) Create(c *gin.Context) {
	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	post := &model.Post{
		Content: content,
		UserID:  currentId.(uint),
	}

	if err := h.svc.Create(c, post); err != nil {
		global.Logger.Errorw("创建动态失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": post})
}

func (h *PostHandler) GetTimeline(c *gin.Context) {
	posts, err := h.svc.GetTimeline(c)
	if err != nil {
		global.Logger.Errorw("获取时间线失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取动态列表失败"})
		return
	}

	if posts == nil {
		posts = []model.Post{}
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func (h *PostHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的动态 ID"})
		return
	}

	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	if err := h.svc.Delete(c, currentId.(uint), uint(id)); err != nil {
		global.Logger.Errorw("删除动态失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
