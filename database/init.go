package database

import (
	"strconv"

	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 初始化数据库
func Init() {
	switch setting.SettingVar.Database.TypeName {
	case "sqlite":
		MainDB, _ = gorm.Open(sqlite.Open(setting.SettingVar.Database.Source))
	case "mysql":
		MainDB, _ = gorm.Open(mysql.Open(setting.SettingVar.Database.Source))
	default:
		panic(setting.SettingVar.Database.TypeName + "是不支持的数据库类型")
	}
	MainDB.AutoMigrate(&UserTable{}, &LinkTable{}, &DataTable{})
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:" + strconv.Itoa(setting.SettingVar.RedisPort),
	})
}
