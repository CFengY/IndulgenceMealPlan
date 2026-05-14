package service

import (
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/middleware"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/repository"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ILoginService interface {
	// Register 注册新用户
	Register(username, password string) (map[string]interface{}, error)

	// Login 验证用户凭据并返回登录响应
	Login(username, password string) (map[string]interface{}, error)

	// Logout 登出用户
	Logout(c *gin.Context) error
}

type LoginService struct {
	repo repository.ILoginRepository
}

func NewLoginService(repo repository.ILoginRepository) ILoginService {
	return &LoginService{repo: repo}
}

func (s *LoginService) Register(username, password string) (map[string]interface{}, error) {
	_, err := s.repo.GetUserByUsername(username)
	if err == nil {
		return nil, fmt.Errorf("用户已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("其他数据库错误: %w", err)
	}

	user := &model.User{
		Username: username,
		Password: password, // In production, hash the password before storing
	}

	s.repo.CreateUser(user)

	token, err := middleware.GenerateToken(uint(user.ID), "签发服务", global.Config.Jwt.Secret)
	if err != nil {
		return nil, fmt.Errorf("用户创建失败: %w", err)
	}

	response := map[string]interface{}{
		"userid":   user.ID,
		"username": username,
		"token":    token,
	}

	return response, nil
}

func (s *LoginService) Login(username, password string) (map[string]interface{}, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("无效的用户名或密码")
	}
	if user.Password != password { // In production, compare hashed passwords
		return nil, fmt.Errorf("无效的用户名或密码")
	}

	token, err := middleware.GenerateToken(uint(user.ID), "签发服务", global.Config.Jwt.Secret)

	if err != nil {
		return nil, fmt.Errorf("token生成失败")
	}

	response := map[string]interface{}{
		"userid":   user.ID,
		"username": username,
		"token":    token,
	}

	return response, nil
}

func (s *LoginService) Logout(c *gin.Context) error {
	token := c.Request.Header.Get(global.Config.Jwt.Name)
	if token == "" {
		return fmt.Errorf("token is required for logout")
	}

	// 处理Bearer前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims := c.MustGet("claims").(*middleware.JwtPayload)

	ctx := c.Request.Context()

	ttl := time.Until(time.Unix(claims.ExpiresAt.Unix(), 0))

	err := global.Redis.Set(ctx, "logout:"+token, "true", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set logout token")
	}

	return nil
}
