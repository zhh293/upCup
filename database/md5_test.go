package database

import (
	"crypto/md5"
	"fmt"
	"io"
	"testing"
)

func TestMD5(t *testing.T) {
	h := md5.New()
	io.WriteString(h, "testpassword")
	res := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(res)
}
