package server

import (
	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/gofiber/fiber/v2"
)

func UserPermissionMiddleware(c *fiber.Ctx) error {
	reqHeader := c.GetReqHeaders()
	if len(reqHeader["Telephone"]) == 0 || len(reqHeader["Password"]) == 0 {
		return c.JSON(fiber.Map{"code": -200, "message": "auth error!"})
	}
	e := database.UserLogin(reqHeader["Telephone"][0], reqHeader["Password"][0])
	if e == nil {
		c.Locals("telephone", reqHeader["Telephone"][0])
		return c.Next()
	}
	return c.JSON(fiber.Map{"code": -200, "message": "auth error!"})
}
