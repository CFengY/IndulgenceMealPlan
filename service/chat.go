package service

import (
	"IndulgenceMealPlan/config"
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

const systemPrompt = `你是一个专业的营养师和饮食顾问，名叫"饮食小助手"。你的职责是根据用户的饮食记录，提供个性化的营养分析和饮食建议。

## 回答规则

1. **基于数据**：引用用户的具体饮食记录进行分析，不要凭空猜测。
2. **专业知识**：你可以回答食物热量、营养成分、膳食搭配等营养学问题。热量估算请标注"大致"。
3. **个性化推荐**：结合用户的饮食偏好和历史记录，推荐合适的食物或餐食搭配。
4. **减脂/增肌建议**：如果用户有健身目标，分析其饮食结构是否合理，给出调整建议。
5. **免责声明**：在涉及健康问题的回复末尾，加上"温馨提示：以上建议仅供参考，不构成专业医疗建议。如有健康问题，请咨询注册营养师或医生。"

## 用户近期饮食记录
{diet_context}

## 回答语言
使用中文回答。保持亲切、鼓励的语气，适当使用 emoji 增加亲和力。`

type ChatService struct {
	mealRepo  repository.IMealRepository
	chatModel *openai.ChatModel
	template  prompt.ChatTemplate
}

func NewChatService(mealRepo repository.IMealRepository, cfg config.AIConfig) (*ChatService, error) {
	ctx := context.Background()

	// fmt.Println(cfg)

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   cfg.Model,
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 ChatModel 失败: %w", err)
	}

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemPrompt),
		schema.UserMessage("{question}"),
	)

	return &ChatService{
		mealRepo:  mealRepo,
		chatModel: chatModel,
		template:  template,
	}, nil
}

func (s *ChatService) Chat(ctx context.Context, userID uint, question string) (<-chan string, error) {
	// 获取用户近 30 天饮食记录
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	meals, err := s.mealRepo.ListByDateRange(userID, startDate, endDate)
	if err != nil {
		global.Logger.Warnw("获取饮食记录失败", "error", err)
	}

	// 格式化为上下文
	dietContext := formatDietContext(meals)

	// fmt.Println(dietContext)

	// 渲染提示词
	messages, err := s.template.Format(ctx, map[string]any{
		"diet_context": dietContext,
		"question":     question,
	})
	if err != nil {
		return nil, fmt.Errorf("提示词渲染失败: %w", err)
	}

	// fmt.Println(s.template)
	// fmt.Println(messages)

	// 调用 LLM
	// result, err := s.chatModel.Generate(ctx, messages)
	results, err := s.chatModel.Stream(ctx, messages)

	// fmt.Println(result)

	if err != nil {
		return nil, fmt.Errorf("AI 调用失败: %w", err)
	}
	// defer results.Close()

	ch := make(chan string)

	go func() {
		defer close(ch)
		defer results.Close()

		for {
			result, err := results.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				global.Logger.Errorw("接收 AI 响应失败", "error", err)
				break
			}

			ch <- result.Content

		}

	}()

	return ch, nil
}

func formatDietContext(meals []model.Meal) string {
	if len(meals) == 0 {
		return "暂无饮食记录。"
	}

	// 按日期分组
	dateMap := make(map[string][]model.Meal)
	for _, meal := range meals {
		dateStr := meal.MealDate.Format("2006-01-02")
		dateMap[dateStr] = append(dateMap[dateStr], meal)
	}

	// 排序日期
	dates := make([]string, 0, len(dateMap))
	for d := range dateMap {
		dates = append(dates, d)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i] > dates[j]
	})

	var sb strings.Builder
	for _, date := range dates {
		sb.WriteString(fmt.Sprintf("- %s：", date))
		foods := dateMap[date]
		for i, m := range foods {
			if i > 0 {
				sb.WriteString("、")
			}
			sb.WriteString(fmt.Sprintf("%s[%s]", m.FoodName, m.MealType))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
