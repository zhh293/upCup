package main

import (
	"github.com/dingdinglz/ai-swindle-detecter-backend/database"
	"github.com/dingdinglz/ai-swindle-detecter-backend/server"
	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
)

func main() {
	setting.SystemPrepare()
	setting.Open()
	database.Init()
	server.Init()
}
