package server

import (
	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/gofiber/fiber/v2"
)

func LinkAddRoute(c *fiber.Ctx) error {
	if c.FormValue("telephone1", "") == "" || c.FormValue("telephone2", "") == "" {
		return c.JSON(fiber.Map{"code": 1, "message": "参数不全"})
	}
	e := database.LinkAdd(c.FormValue("telephone1"), c.FormValue("telephone2"))
	if e != nil {
		return c.JSON(fiber.Map{"code": 2, "message": "关联已存在"})
	}
	return c.JSON(fiber.Map{"code": 0, "message": ""})
}

func LinkExsitRoute(c *fiber.Ctx) error {
	if c.FormValue("telephone1", "") == "" || c.FormValue("telephone2", "") == "" {
		return c.JSON(fiber.Map{"code": 1, "message": "参数不全"})
	}
	if database.LinkExsit(c.FormValue("telephone1"), c.FormValue("telephone2")) {
		return c.JSON(fiber.Map{"code": 0, "message": "", "exist": true})
	}
	return c.JSON(fiber.Map{"code": 0, "message": "", "exist": false})
}

func LinkGetRoute(c *fiber.Ctx) error {
	if c.FormValue("telephone", "") == "" {
		return c.JSON(fiber.Map{"code": 1, "message": "参数不全"})
	}
	req := c.GetReqHeaders()
	if req["Telephone"][0] != c.FormValue("telephone") {
		return c.JSON(fiber.Map{"code": 2, "message": "权限错误"})
	}
	res := database.LinkGetAll(c.FormValue("telephone"))
	return c.JSON(fiber.Map{"code": 0, "message": "", "data": res})
}
