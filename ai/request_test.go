package ai

import (
	"fmt"
	"testing"
)

func ToolTestRun(sentence string, t *testing.T) {
	res := Run(sentence, 6666)
	if res == "err" {
		t.Error("ai run error!")
		return
	}
	fmt.Println(res)
}
func TestRun(t *testing.T) {
	ToolTestRun("不交保证金，不交会费，即可赚取零花钱，最适合宝妈和学生。", t)
	ToolTestRun("教大家一个网上日赚千元的方法，手机在家就可以做的兼职。", t)
	ToolTestRun("有一笔100万的扶贫金，按流程操作就可以免费申领！", t)
	ToolTestRun("在么?最近出了点事儿，急需用钱，给哥们儿借点", t)
	ToolTestRun("下午好，吃饭了吗", t)
}
