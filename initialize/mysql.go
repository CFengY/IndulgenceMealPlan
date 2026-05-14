package initialize

import (
	"IndulgenceMealPlan/global"
	"IndulgenceMealPlan/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL() {
	// 连接数据库
	dsn := global.Config.MySQL.Username + ":" + global.Config.MySQL.Password + "@tcp(" + global.Config.MySQL.Host + ":" + fmt.Sprint(global.Config.MySQL.Port) + ")/" + global.Config.MySQL.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		global.Logger.Fatalw("连接数据库失败", "error", err)
	}

	// 自动迁移模型
	if err := db.AutoMigrate(&model.User{}, &model.Meal{}); err != nil {
		global.Logger.Fatalw("自动迁移失败", "error", err)
	}

	global.DB = db
}

func CloseMySQL() {
	if global.DB != nil {
		sqlDB, err := global.DB.DB()
		if err != nil {
			global.Logger.Errorw("获取底层sql.DB失败", "error", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			global.Logger.Errorw("关闭MySQL连接失败", "error", err)
			return
		}
		global.Logger.Info("MySQL连接已关闭")
	}
}
