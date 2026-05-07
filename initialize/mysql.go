package initialize

import (
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL() {
	// 连接数据库
	dsn := global.Config.MySQL.Username + ":" + global.Config.MySQL.Password + "@tcp(" + global.Config.MySQL.Host + ":" + fmt.Sprint(global.Config.MySQL.Port) + ")/" + global.Config.MySQL.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移模型
	if err := db.AutoMigrate(&model.User{}, &model.Meal{}); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	global.DB = db
}
