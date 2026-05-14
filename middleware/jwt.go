package middleware

import (
	"IndulgenceMealPlan/global"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	jwt.RegisteredClaims
	UserId uint
}

func GenerateToken(uid uint, subject string, secret string) (string, error) {
	claim := JwtPayload{
		UserId: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                 //签发者
			Subject:   subject,                                       //签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},      //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), //过期时间
			// NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt: jwt.NewNumericDate(time.Now()), //签发时间
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
	return token, err
}

func ParseToken(token string, secret string) (*JwtPayload, error) {
	// 解析token
	parseToken, err := jwt.ParseWithClaims(token, &JwtPayload{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parseToken.Claims.(*JwtPayload); ok && parseToken.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func VerifyJWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(global.Config.Jwt.Name)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			c.Abort()
			return
		}

		// 处理Bearer前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		if global.Redis.Get(c.Request.Context(), "logout:"+token).Val() == "true" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			c.Abort()
			return
		}

		// 解析获取用户载荷信息
		payLoad, err := ParseToken(token, global.Config.Jwt.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			c.Abort()
			return
		}
		// 在上下文设置载荷信息
		c.Set("currentId", payLoad.UserId)
		c.Set("claims", payLoad)
		// 这里是否要通知客户端重新保存新的Token
		c.Next()
	}
}
