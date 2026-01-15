package server

import (
	"context"
	"strconv"

	"github.com/dingdinglz/ai-swindle-detecter-backend/ai"
	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
	"github.com/gofiber/fiber/v2"
)

func AIApiRoute(c *fiber.Ctx) error {
	if c.FormValue("text", "") == "" {
		return c.JSON(fiber.Map{"code": -1, "message": "参数不全"})
	}
	telephone := c.Locals("telephone").(string)
	redisContext := context.Background()
	if setting.SettingVar.Debug {
		_, e := database.RedisClient.Get(redisContext, telephone+":0").Result()
		if e != nil {
			database.RedisClient.Set(redisContext, telephone+":0", int64(1), 0)
		} else {
			database.RedisClient.Incr(redisContext, telephone+":0")
		}
		return c.JSON(fiber.Map{"code": 0, "message": "", "type": "中性"})
	}
	res := ai.Run(c.FormValue("text", ""), setting.SettingVar.AIPort)
	if res == "err" {
		return c.JSON(fiber.Map{"code": 1, "message": "ai error"})
	}
	resMap := make(map[string]int)
	resMap["中性"] = 0
	resMap["网络交易及兼职诈骗"] = 1
	resMap["虚假金融及投资诈骗"] = 2
	resMap["身份冒充及威胁诈骗"] = 3
	_, e := database.RedisClient.Get(redisContext, telephone+":"+strconv.Itoa(resMap[res])).Result()
	if e != nil {
		database.RedisClient.Set(redisContext, telephone+":"+strconv.Itoa(resMap[res]), int64(1), 0)
	} else {
		database.RedisClient.Incr(redisContext, telephone+":"+strconv.Itoa(resMap[res]))
	}
	return c.JSON(fiber.Map{"code": 0, "message": "", "type_id": resMap[res], "type": res})
}
