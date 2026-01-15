package server

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSlice(t *testing.T) {
	var text []string
	for i := 0; i <= 100; i++ {
		text = append(text, strconv.Itoa(i))
	}
	fmt.Println(text[5])
	fmt.Println(text[0:5])
	fmt.Println(text[99:101])
}
