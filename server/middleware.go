package server

import (
	"strings"

	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
	"github.com/dingdinglz/ai-swindle-detecter-backend/tools"
	"github.com/gofiber/fiber/v2"
)

func UserPermissionMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization", "")
	if authHeader == "" {
		return c.JSON(fiber.Map{"code": "401", "msg": "未提供认证token"})
	}

	// 解析Authorization头部，格式为 "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.JSON(fiber.Map{"code": "401", "msg": "token格式错误"})
	}

	tokenString := parts[1]
	claims, err := tools.ParseToken(tokenString, setting.SettingVar.JWT.Secret)
	if err != nil {
		return c.JSON(fiber.Map{"code": "401", "msg": "无效的token"})
	}

	// 将用户信息存储到上下文中
	c.Locals("telephone", claims.Telephone)
	c.Locals("user_id", claims.UserID)

	return c.Next()
}

