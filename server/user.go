package server

import (
	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/gofiber/fiber/v2"
)

func UserRegisterRoute(c *fiber.Ctx) error {
	if c.FormValue("telephone", "") == "" || c.FormValue("password", "") == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	e := database.UserNew(c.FormValue("telephone"), c.FormValue("password"))
	if e != nil {
		return c.JSON(fiber.Map{"code": 1, "message": e.Error()})
	}
	return c.JSON(fiber.Map{"code": 0, "message": ""})
}

func UserLoginRoute(c *fiber.Ctx) error {
	if c.FormValue("telephone", "") == "" || c.FormValue("password", "") == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	e := database.UserLogin(c.FormValue("telephone"), c.FormValue("password"))
	if e == nil {
		return c.JSON(fiber.Map{"code": 0, "message": ""})
	}
	if e.Error() == "not exist" {
		return c.JSON(fiber.Map{"code": 1, "message": "用户不存在"})
	}
	return c.JSON(fiber.Map{"code": 2, "message": "密码错误"})
}
