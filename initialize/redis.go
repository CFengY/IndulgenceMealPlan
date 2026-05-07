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
		panic(err)
	}
	global.Redis = client
}
