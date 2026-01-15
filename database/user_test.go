package database

import (
	"testing"

	"github.com/dingdinglz/ai-swindle-detecter-backend/setting"
)

func TestUserCreate(t *testing.T) {
	setting.RootPath = "../"
	Init()
	UserNew("15255601211", "testpassword")
}
