package initialize

import (
	"IndulgenceMealPlan/global"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	redisOpt := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisOpt.Host, redisOpt.Port),
		Password: redisOpt.Password, // no password set
		DB:       redisOpt.DataBase, // use default DB
	})
	ping := client.Ping(context.Background())
	err := ping.Err()
	if err != nil {
		global.Logger.Fatalw("Redis连接失败", "error", err)
	}
	global.Redis = client
}

func CloseRedis() {
	if global.Redis != nil {
		if err := global.Redis.Close(); err != nil {
			global.Logger.Errorw("关闭Redis连接失败", "error", err)
			return
		}
		global.Logger.Info("Redis连接已关闭")
	}
}
