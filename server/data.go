package server

import (
	"context"
	"strconv"

	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/dingdinglz/ai-swindle-detecter-backend/tools"
	"github.com/gofiber/fiber/v2"
)

func DataAddRoute(c *fiber.Ctx) error {
	var req struct {
		Telephone string `json:"telephone" form:"telephone"`
		Type      string `json:"type" form:"type"`
		Text      string `json:"text" form:"text"`
		Package   string `json:"package" form:"package"`
	}
	if err := c.BodyParser(&req); err != nil {
		// 忽略解析错误
	}
	
	// 回退机制：如果 BodyParser 没解析到，尝试 FormValue
	if req.Telephone == "" { req.Telephone = c.FormValue("telephone") }
	if req.Type == "" { req.Type = c.FormValue("type") }
	if req.Text == "" { req.Text = c.FormValue("text") }
	if req.Package == "" { req.Package = c.FormValue("package") }

	if req.Telephone == "" || req.Type == "" || req.Text == "" || req.Package == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	// 从 Token 获取当前操作用户的手机号
	currentUser := c.Locals("telephone").(string)
	targetUser := req.Telephone

	// 检查权限：操作人必须是本人或者与目标用户存在关联
	if !database.UserCheckAllow(targetUser, currentUser) {
		return c.JSON(fiber.Map{"code": 1, "message": "权限错误"})
	}
	database.DataAdd(req.Package, targetUser, req.Text, req.Type)
	return c.JSON(fiber.Map{"code": 0, "message": ""})
}

func DataGetRoute(c *fiber.Ctx) error {
	if c.Query("telephone", "") == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	currentUser := c.Locals("telephone").(string)
	targetUser := c.Query("telephone", "")

	if !database.UserCheckAllow(targetUser, currentUser) {
		return c.JSON(fiber.Map{"code": 1, "message": "权限错误"})
	}
	i := database.DataGet(targetUser)
	return c.JSON(fiber.Map{"code": 0, "message": "", "data": i})
}

func DataCutGetRoute(c *fiber.Ctx) error {
	if c.FormValue("page", "") == "" || c.FormValue("telephone", "") == "" || c.FormValue("cut", "") == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	currentUser := c.Locals("telephone").(string)
	targetUser := c.FormValue("telephone", "")

	if !database.UserCheckAllow(targetUser, currentUser) {
		return c.JSON(fiber.Map{"code": 1, "message": "权限错误"})
	}
	page := tools.StringToInt(c.FormValue("page"))
	cut := tools.StringToInt(c.FormValue("cut"))
	counts := database.DataCounts(c.FormValue("telephone"))
	if counts == 0 {
		return c.JSON(fiber.Map{"code": 0, "message": "", "pages": 0, "data": []database.DataTable{}})
	}
	pages := counts / cut
	if counts%cut != 0 {
		pages++
	}
	if page <= 0 || page > pages {
		return c.JSON(fiber.Map{"code": 2, "message": "参数错误"})
	}
	datas := database.DataGet(c.FormValue("telephone"))
	start := (page - 1) * cut
	end := start + cut
	if end > counts {
		end = counts
	}
	return c.JSON(fiber.Map{"code": 0, "message": "", "pages": pages, "data": datas[start:end]})
}

func DataCountRoute(ctx *fiber.Ctx) error {
	telephone := ctx.Locals("telephone").(string)
	var all int64 = 0
	datas := make(map[int]int64)
	for i := 0; i <= 3; i++ {
		data, e := database.RedisClient.Get(context.Background(), telephone+":"+strconv.Itoa(i)).Int64()
		if e != nil {
			data = 0
		}
		datas[i] = data
		all += data
	}
	return ctx.JSON(fiber.Map{
		"code": 0,
		"data": datas,
		"all":  all,
	})
}
