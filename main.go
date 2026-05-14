package main

import (
	"IndulgenceMealPlan/config"
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/handler"
	"IndulgenceMealPlan/initialize"
	"IndulgenceMealPlan/repository"
	"IndulgenceMealPlan/router"
	"IndulgenceMealPlan/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	global.Config, err = config.Load()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	initialize.InitLogger()
	defer global.Logger.Sync()

	// 初始化上传目录
	if err := os.MkdirAll(global.Config.Upload.Dir, 0755); err != nil {
		global.Logger.Fatalw("创建上传目录失败", "error", err)
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
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	router.Setup(r, mealHandler, loginHandler, global.Config.Upload.Dir)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", global.Config.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 在 goroutine 中启动服务
	go func() {
		global.Logger.Infow("服务启动", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Fatalw("服务启动失败", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	global.Logger.Infow("收到关闭信号", "signal", sig.String())

	// 优雅关闭
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(global.Config.Server.ShutdownTimeout)*time.Second,
	)
	defer cancel()

	global.Logger.Info("正在关闭HTTP服务...")
	if err := srv.Shutdown(ctx); err != nil {
		global.Logger.Errorw("HTTP服务强制关闭", "error", err)
	}

	// 关闭数据库连接
	initialize.CloseMySQL()
	initialize.CloseRedis()

	global.Logger.Info("服务已关闭")
}
