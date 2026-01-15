package server

import (
	"context"
	"strconv"

	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/dingdinglz/ai-swindle-detecter-backend/tools"
	"github.com/gofiber/fiber/v2"
)

func DataAddRoute(c *fiber.Ctx) error {
	if c.FormValue("telephone", "") == "" || c.FormValue("type", "") == "" || c.FormValue("text", "") == "" || c.FormValue("package", "") == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	headMap := c.GetReqHeaders()
	if !database.UserCheckAllow(c.FormValue("telephone", ""), headMap["Telephone"][0]) {
		return c.JSON(fiber.Map{"code": 1, "message": "权限错误"})
	}
	database.DataAdd(c.FormValue("package"), c.FormValue("telephone"), c.FormValue("text"), c.FormValue("type"))
	return c.JSON(fiber.Map{"code": 0, "message": ""})
}

func DataGetRoute(c *fiber.Ctx) error {
	if c.Query("telephone", "") == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	headMap := c.GetReqHeaders()
	if !database.UserCheckAllow(c.Query("telephone"), headMap["Telephone"][0]) {
		return c.JSON(fiber.Map{"code": 1, "message": "权限错误"})
	}
	i := database.DataGet(c.Query("telephone"))
	return c.JSON(fiber.Map{"code": 0, "message": "", "data": i})
}

func DataCutGetRoute(c *fiber.Ctx) error {
	if c.FormValue("page", "") == "" || c.FormValue("telephone", "") == "" || c.FormValue("cut", "") == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	headMap := c.GetReqHeaders()
	if !database.UserCheckAllow(c.FormValue("telephone"), headMap["Telephone"][0]) {
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
