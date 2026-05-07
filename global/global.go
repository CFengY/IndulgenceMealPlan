package global

import (
	"IndulgenceMealPlan/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
)
