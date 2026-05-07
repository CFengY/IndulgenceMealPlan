package main

import (
	"IndulgenceMealPlan/config"
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/handler"
	"IndulgenceMealPlan/initialize"
	"IndulgenceMealPlan/repository"
	"IndulgenceMealPlan/router"
	"IndulgenceMealPlan/service"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	global.Config, err = config.Load()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化上传目录
	if err := os.MkdirAll(global.Config.Upload.Dir, 0755); err != nil {
		fmt.Printf("创建上传目录失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库连接
	initialize.InitMySQL()
	initialize.InitRedis()

	// 初始化各层
	mealRepo := repository.NewMealRepository(global.DB)
	mealSvc := service.NewMealService(mealRepo, global.Config.Upload.Dir, global.Config.Upload.MaxSize)
	mealHandler := handler.NewMealHandler(mealSvc)

	loginRepo := repository.NewLoginRepository(global.DB)
	loginSvc := service.NewLoginService(loginRepo)
	loginHandler := handler.NewLoginHandler(loginSvc)

	// 初始化路由
	r := gin.Default()
	router.Setup(r, mealHandler, loginHandler, global.Config.Upload.Dir)

	// 启动服务
	addr := fmt.Sprintf(":%d", global.Config.Server.Port)
	fmt.Printf("服务启动在 %s\n", addr)
	if err := r.Run(addr); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
		os.Exit(1)
	}
}
