package handler

import (
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	svc *service.ChatService
}

func NewChatHandler(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入消息"})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息不能为空"})
		return
	}

	currentId, exists := c.Get("currentId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	reply, err := h.svc.Chat(c.Request.Context(), currentId.(uint), req.Message)
	if err != nil {
		global.Logger.Errorw("AI 聊天失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 服务暂时不可用，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, model.ChatResponse{Reply: reply})
}
