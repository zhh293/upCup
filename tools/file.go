package tools

import "os"

// 判断文件或文件夹是否存在
func IsFileOrDirExist(_path string) bool {
	_, e := os.Stat(_path)
	return e == nil
}

// 若文件夹不存在则创建文件夹
func MkdirINE(_path string) {
	if !IsFileOrDirExist(_path) {
		os.Mkdir(_path, os.ModePerm)
	}
}
