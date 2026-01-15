package database

import (
	"testing"

	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
)

func TestLink(t *testing.T) {
	setting.RootPath = "../"
	Init()
	t.Log(LinkExsit("15255601211", "15255601212"))
}
