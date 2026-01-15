package server

import (
	"fmt"

	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// 初始化并启动服务器
func Init() {
	MainServer = fiber.New()
	MainServer.Use(logger.New(), recover.New(), cors.New())
	MainServer.Get("/monitor", monitor.New())
	BindRoutes()
	e := MainServer.Listen("0.0.0.0:" + setting.SettingVar.Port)
	if e != nil {
		fmt.Println(e.Error())
	}
}
