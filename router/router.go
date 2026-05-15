package router

import (
	"IndulgenceMealPlan/handler"
	"IndulgenceMealPlan/middleware"

	corspkg "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, hMeal *handler.MealHandler, hLogin *handler.LoginHandler, hChat *handler.ChatHandler, uploadDir string) {
	r.Use(middleware.RequestLogger())
	r.Use(corspkg.New(middleware.CORS()))

	api := r.Group("/api/v1")
	api.Use(middleware.RateLimitMiddleware())
	api.Use(middleware.VerifyJWTAdmin())
	{
		api.POST("/meals", hMeal.Create)
		api.GET("/meals/range", hMeal.ListByDateRange)
		api.GET("/meals", hMeal.GetMeals)
		api.PUT("/meals/:id", hMeal.Update)
		api.DELETE("/meals/:id", hMeal.Delete)

		api.POST("/logout", hLogin.Logout)
		api.POST("/chat", hChat.Chat)
	}

	authApi := r.Group("/api/v1/auth")
	{
		authApi.POST("/login", hLogin.Login)
		authApi.POST("/register", hLogin.Register)
	}

	r.Static("/images", uploadDir)
}
