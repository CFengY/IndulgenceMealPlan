package handler

import (
	"IndulgenceMealPlan/middleware"
	"IndulgenceMealPlan/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Register(username, password string) (map[string]interface{}, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockLoginService) Login(username, password string) (map[string]interface{}, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockLoginService) Logout(c *gin.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestLoginHandler_Register_Success(t *testing.T) {
	mockSvc := new(MockLoginService)
	response := map[string]interface{}{
		"userid":   float64(1),
		"username": "testuser",
		"token":    "test-token",
	}
	mockSvc.On("Register", "testuser", "password123").Return(response, nil)
	handler := NewLoginHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/auth/register", url.Values{
		"username": {"testuser"},
		"password": {"password123"},
	})

	handler.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLoginHandler_Register_MissingFields(t *testing.T) {
	mockSvc := new(MockLoginService)
	handler := NewLoginHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/auth/register", url.Values{})

	handler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler_Login_Success(t *testing.T) {
	mockSvc := new(MockLoginService)
	response := map[string]interface{}{
		"userid":   float64(1),
		"username": "testuser",
		"token":    "test-token",
	}
	mockSvc.On("Login", "testuser", "password123").Return(response, nil)
	handler := NewLoginHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/auth/login", url.Values{
		"username": {"testuser"},
		"password": {"password123"},
	})

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLoginHandler_Login_InvalidCredentials(t *testing.T) {
	mockSvc := new(MockLoginService)
	mockSvc.On("Login", "testuser", "wrongpass").Return(nil, errors.New("无效的用户名或密码"))
	handler := NewLoginHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/auth/login", url.Values{
		"username": {"testuser"},
		"password": {"wrongpass"},
	})

	handler.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginHandler_Logout_Success(t *testing.T) {
	mockSvc := new(MockLoginService)
	mockSvc.On("Logout", mock.Anything).Return(nil)
	handler := NewLoginHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")
	// 模拟 JWT 中间件设置的 claims
	c.Set("claims", &middleware.JwtPayload{})

	handler.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLoginHandler_Logout_NoToken(t *testing.T) {
	mockSvc := new(MockLoginService)
	handler := NewLoginHandler(mockSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)

	handler.Logout(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler_Register_DuplicateUser(t *testing.T) {
	mockSvc := new(MockLoginService)
	mockSvc.On("Register", "existing", "pass123").Return(nil, errors.New("用户已存在"))
	handler := NewLoginHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/auth/register", url.Values{
		"username": {"existing"},
		"password": {"pass123"},
	})

	handler.Register(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLoginHandler_Login_MissingFields(t *testing.T) {
	mockSvc := new(MockLoginService)
	handler := NewLoginHandler(mockSvc)

	c, w := setupTestContext(http.MethodPost, "/api/v1/auth/login", url.Values{
		"username": {""},
		"password": {""},
	})

	handler.Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Register 返回注册响应
func (*MockLoginService) unused() {
	_ = service.ILoginService(nil)
}

// 确保 MockLoginService 实现了 ILoginService 接口
var _ service.ILoginService = (*MockLoginService)(nil)

func TestMockLoginService_ImplementsInterface(t *testing.T) {
	// 编译时已检查，此处仅确保 mock 可用
	svc := new(MockLoginService)
	assert.Implements(t, (*service.ILoginService)(nil), svc)
}
