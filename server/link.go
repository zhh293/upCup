package server

import (
	"strconv"

	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/gofiber/fiber/v2"
)

type linkAddRequest struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

func LinkListRoute(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
	}

	linkedAccounts, err := database.LinkedAccountList(userID)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    500,
			"message": "获取关联账户失败",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    linkedAccounts,
	})
}

func LinkAddRoute(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	userID, ok := userIDValue.(string)
	telephoneValue := c.Locals("telephone")
	userTelephone, ok2 := telephoneValue.(string)

	if !ok || userID == "" || !ok2 || userTelephone == "" {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
	}

	var req linkAddRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "请求体解析失败",
			"data":    nil,
		})
	}

	if req.Name == "" || req.Telephone == "" {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "参数不全",
			"data":    nil,
		})
	}

	if req.Telephone == userTelephone {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "不能关联自己",
			"data":    nil,
		})
	}

	record, err := database.LinkedAccountAdd(userID, userTelephone, req.Name, req.Telephone)
	if err != nil {
		if err.Error() == "account already linked" {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "该账号已关联",
				"data":    nil,
			})
		}
		if err.Error() == "target user not found" {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "目标用户不存在",
				"data":    nil,
			})
		}
		return c.JSON(fiber.Map{
			"code":    500,
			"message": "添加失败",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "添加成功",
		"data":    record,
	})
}

func LinkRemoveRoute(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	userID, ok := userIDValue.(string)
	telephoneValue := c.Locals("telephone")
	userTelephone, ok2 := telephoneValue.(string)

	if !ok || userID == "" || !ok2 || userTelephone == "" {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "未认证",
			"data":    nil,
		})
	}

	idParam := c.Params("id", "")
	if idParam == "" {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "参数不全",
			"data":    nil,
		})
	}

	idValue, err := strconv.Atoi(idParam)
	if err != nil || idValue <= 0 {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
	}

	err = database.LinkedAccountRemove(userID, userTelephone, uint(idValue))
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "关联账户不存在",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "移除成功",
		"data":    nil,
	})
}

func LinkGetSMSRoute(c *fiber.Ctx) error {
	currentUser := c.Locals("telephone").(string)
	targetUser := c.Query("telephone", "")

	if targetUser == "" {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "参数不全",
			"data":    nil,
		})
	}

	// 检查权限 (UserCheckAllow 会检查是否是本人或已存在关联)
	if !database.UserCheckAllow(targetUser, currentUser) {
		return c.JSON(fiber.Map{
			"code":    403,
			"message": "权限错误或未关联",
			"data":    nil,
		})
	}

	data := database.DataGet(targetUser)
	return c.JSON(fiber.Map{
		"code":    0,
		"message": "",
		"data":    data,
	})
}
