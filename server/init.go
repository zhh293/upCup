package server

import (
	"fmt"
	"path/filepath"

	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// 初始化并启动服务器
func Init() {
	MainServer = fiber.New(fiber.Config{
		BodyLimit: 55 * 1024 * 1024, // 55MB
	})
	MainServer.Use(logger.New(), recover.New(), cors.New())
	MainServer.Get("/monitor", monitor.New())
	MainServer.Static("/static/audio", filepath.Join(setting.RootPath, "data", "audio"))
	BindRoutes()
	e := MainServer.Listen("0.0.0.0:" + setting.SettingVar.Port)
	if e != nil {
		fmt.Println(e.Error())
	}
}
