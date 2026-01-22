package setting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dingdinglz/ai-swindle-detecter-backend/tools"
)

var (
	SettingVar SettingModel // 设置内容
	RootPath   string       // 运行地址
)

// 创建数据文件夹等
func SystemPrepare() {
	RootPath, _ = os.Getwd()
	tools.MkdirINE(filepath.Join(RootPath, "data"))
}

// 打开设置
func Open() {
	if !tools.IsFileOrDirExist(filepath.Join(RootPath, "data", "setting.json")) {
		SettingVar.Port = "7000"
		SettingVar.AIHost = "127.0.0.1"
		SettingVar.AIPort = 6666
		SettingVar.Debug = false
		SettingVar.Database.TypeName = "sqlite"
		SettingVar.Database.Source = filepath.Join(RootPath, "data", "data.db")
		SettingVar.RedisPort = 6379
		SettingVar.JWT.Secret = "your-secret-key-change-this-in-production"
		SettingVar.JWT.ExpiresIn = 7200
		SettingVar.JWT.TokenType = "Bearer"
		r, _ := json.Marshal(SettingVar)
		os.WriteFile(filepath.Join(RootPath, "data", "setting.json"), r, os.ModePerm)
	} else {
		j, _ := os.ReadFile(filepath.Join(RootPath, "data", "setting.json"))
		json.Unmarshal(j, &SettingVar)
		// 确保JWT配置有默认值
		if SettingVar.JWT.Secret == "" {
			SettingVar.JWT.Secret = "your-secret-key-change-this-in-production"
		}
		if SettingVar.JWT.ExpiresIn == 0 {
			SettingVar.JWT.ExpiresIn = 7200
		}
		if SettingVar.JWT.TokenType == "" {
			SettingVar.JWT.TokenType = "Bearer"
		}
		if SettingVar.AIHost == "" {
			SettingVar.AIHost = "127.0.0.1"
		}
	}
	if SettingVar.Debug {
		fmt.Println("[Debug Mode]该模式下AIapi进行假调用，返回一个随机的类型")
	}
}
