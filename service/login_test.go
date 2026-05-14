package service

import (
	"IndulgenceMealPlan/config"
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"IndulgenceMealPlan/repository"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupLoginTest() {
	if global.Config == nil {
		global.Config = &config.Config{
			Jwt: config.JwtConfig{
				Secret:     "test-secret",
				Expiration: 3600,
				Name:       "Authorization",
			},
		}
	}
	if global.Redis == nil {
		mr, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		global.Redis = redis.NewClient(&redis.Options{Addr: mr.Addr()})
	}
	gin.SetMode(gin.TestMode)
}

func TestLoginService_Register_Success(t *testing.T) {
	setupLoginTest()

	mockRepo := new(repository.MockLoginRepository)
	mockRepo.On("GetUserByUsername", "newuser").Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)

	svc := NewLoginService(mockRepo)

	result, err := svc.Register("newuser", "password123")
	require.NoError(t, err)
	assert.Equal(t, "newuser", result["username"])
	assert.NotEmpty(t, result["token"])
	mockRepo.AssertExpectations(t)
}

func TestLoginService_Register_DuplicateUser(t *testing.T) {
	setupLoginTest()

	mockRepo := new(repository.MockLoginRepository)
	existingUser := &model.User{Username: "existing"}
	existingUser.ID = 1
	mockRepo.On("GetUserByUsername", "existing").Return(existingUser, nil)

	svc := NewLoginService(mockRepo)

	_, err := svc.Register("existing", "password123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户已存在")
}

func TestLoginService_Login_Success(t *testing.T) {
	setupLoginTest()

	mockRepo := new(repository.MockLoginRepository)
	user := &model.User{Username: "testuser", Password: "password123"}
	user.ID = 1
	mockRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	svc := NewLoginService(mockRepo)

	result, err := svc.Login("testuser", "password123")
	require.NoError(t, err)
	assert.Equal(t, "testuser", result["username"])
	assert.NotEmpty(t, result["token"])
}

func TestLoginService_Login_WrongPassword(t *testing.T) {
	setupLoginTest()

	mockRepo := new(repository.MockLoginRepository)
	user := &model.User{Username: "testuser", Password: "correctpass"}
	user.ID = 1
	mockRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	svc := NewLoginService(mockRepo)

	_, err := svc.Login("testuser", "wrongpass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的用户名或密码")
}

func TestLoginService_Login_UserNotFound(t *testing.T) {
	setupLoginTest()

	mockRepo := new(repository.MockLoginRepository)
	mockRepo.On("GetUserByUsername", "nonexistent").Return(nil, errors.New("not found"))

	svc := NewLoginService(mockRepo)

	_, err := svc.Login("nonexistent", "pass")
	assert.Error(t, err)
}
