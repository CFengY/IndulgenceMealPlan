package middleware

import (
	"IndulgenceMealPlan/config"
	"IndulgenceMealPlan/global"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	global.Config = &config.Config{
		Jwt: config.JwtConfig{
			Secret:     "test-secret-key",
			Expiration: 3600,
			Name:       "Authorization",
		},
	}

	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer mr.Close()
	global.Redis = redis.NewClient(&redis.Options{Addr: mr.Addr()})

	gin.SetMode(gin.TestMode)

	code := m.Run()
	os.Exit(code)
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(1, "test-subject", "my-secret")

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// 解析验证
	claims, err := ParseToken(token, "my-secret")
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserId)
	assert.Equal(t, "test-subject", claims.Subject)
}

func TestParseToken_ValidToken(t *testing.T) {
	token, err := GenerateToken(42, "user42", "secret123")
	require.NoError(t, err)

	claims, err := ParseToken(token, "secret123")
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserId)
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken(1, "test", "correct-secret")
	require.NoError(t, err)

	_, err = ParseToken(token, "wrong-secret")
	assert.Error(t, err)
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := ParseToken("invalid.token.here", "secret")
	assert.Error(t, err)
}

func TestVerifyJWTAdmin_NoToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/meals", nil)

	handler := VerifyJWTAdmin()
	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted())
}

func TestVerifyJWTAdmin_ValidToken(t *testing.T) {
	token, err := GenerateToken(1, "test", global.Config.Jwt.Secret)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/meals", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := VerifyJWTAdmin()
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	currentId, exists := c.Get("currentId")
	assert.True(t, exists)
	assert.Equal(t, uint(1), currentId)
}
