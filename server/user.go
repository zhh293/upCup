package server

import (
	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
	"github.com/dingdinglz/ai-swindle-detecter-backend/tools"
	"github.com/gofiber/fiber/v2"
)

func UserRegisterRoute(c *fiber.Ctx) error {
	telephone := c.FormValue("telephone", "")
	password := c.FormValue("password", "")

	if telephone == "" || password == "" {
		return c.JSON(fiber.Map{"code": "400", "msg": "参数不全"})
	}

	userInfo, err := database.UserNew(telephone, password)
	if err != nil {
		if err.Error() == "user has exists" {
			return c.JSON(fiber.Map{"code": "400", "msg": "用户已存在"})
		}
		return c.JSON(fiber.Map{"code": "500", "msg": "注册失败"})
	}

	// 生成token
	token, err := tools.GenerateToken(telephone, userInfo.UserID, setting.SettingVar.JWT.Secret, setting.SettingVar.JWT.ExpiresIn)
	if err != nil {
		return c.JSON(fiber.Map{"code": "500", "msg": "生成token失败"})
	}

	return c.JSON(fiber.Map{
		"code": "200",
		"msg":  "注册成功",
		"data": fiber.Map{
			"user_id":      userInfo.UserID,
			"telephone":    userInfo.Telephone,
			"access_token": token,
			"token_type":  setting.SettingVar.JWT.TokenType,
			"expires_in":   setting.SettingVar.JWT.ExpiresIn,
		},
	})
}

func UserLoginRoute(c *fiber.Ctx) error {
	telephone := c.FormValue("telephone", "")
	password := c.FormValue("password", "")

	if telephone == "" || password == "" {
		return c.JSON(fiber.Map{"code": "400", "msg": "参数不全"})
	}

	userInfo, err := database.UserLogin(telephone, password)
	if err != nil {
		if err.Error() == "not exist" {
			return c.JSON(fiber.Map{"code": "401", "msg": "用户不存在"})
		}
		return c.JSON(fiber.Map{"code": "401", "msg": "密码错误"})
	}

	// 生成token
	token, err := tools.GenerateToken(telephone, userInfo.UserID, setting.SettingVar.JWT.Secret, setting.SettingVar.JWT.ExpiresIn)
	if err != nil {
		return c.JSON(fiber.Map{"code": "500", "msg": "生成token失败"})
	}

	return c.JSON(fiber.Map{
		"code": "200",
		"msg":  "登录成功",
		"data": fiber.Map{
			"user_id":      userInfo.UserID,
			"telephone":    userInfo.Telephone,
			"nickname":     userInfo.Nickname,
			"avatar":       userInfo.Avatar,
			"email":        userInfo.Email,
			"access_token": token,
			"token_type":   setting.SettingVar.JWT.TokenType,
			"expires_in":   setting.SettingVar.JWT.ExpiresIn,
		},
	})
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func UserLogoutRoute(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "登出成功",
		"data":    nil,
	})
}

func UserChangePasswordRoute(c *fiber.Ctx) error {
	telephoneValue := c.Locals("telephone")
	telephone, ok := telephoneValue.(string)
	if !ok || telephone == "" {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
	}

	var req changePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "请求体解析失败",
			"data":    nil,
		})
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "参数不全",
			"data":    nil,
		})
	}

	err := database.ChangePassword(telephone, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "password incorrect" {
			return c.JSON(fiber.Map{
				"code":    401,
				"message": "原密码错误",
				"data":    nil,
			})
		}
		return c.JSON(fiber.Map{
			"code":    500,
			"message": "密码修改失败",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "密码修改成功，请重新登录",
		"data":    nil,
	})
}

