package handler

import (
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/service"
	"encoding/json"
	"fmt"
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

	ch, err := h.svc.Chat(c.Request.Context(), currentId.(uint), req.Message)
	if err != nil {
		global.Logger.Errorw("AI 聊天失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 服务暂时不可用，请稍后再试"})
		return
	}

	// SSE 流式推送
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		global.Logger.Error("ResponseWriter 不支持 Flusher 接口")
		return
	}

	for content := range ch {
		select {
		case <-c.Request.Context().Done():
			return
		default:
		}

		data, _ := json.Marshal(map[string]string{"content": content})
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
	}

	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	flusher.Flush()
}
