package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func IndexRoute(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"code": 0, "message": "running...", "time": time.Now().Format("2006-01-02 15:04:05.99")})
}
