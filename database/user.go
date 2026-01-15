package database

import (
	"errors"

	"github.com/dingdinglz/ai-swindle-detecter-backend/tools"
)

func UserNew(telephone string, password string) error {
	if UserExist(telephone) {
		return errors.New("user has exists")
	}
	var i UserTable
	i.Telephone = telephone
	i.Password = tools.MD5(password)
	MainDB.Create(&i)
	return nil
}

func UserExist(telephone string) bool {
	var i int64
	MainDB.Model(&UserTable{}).Where("telephone = ?", telephone).Count(&i)
	return i != 0
}

func UserLogin(telephone string, password string) error {
	if !UserExist(telephone) {
		return errors.New("not exist")
	}
	var i UserTable
	MainDB.Model(&UserTable{}).Where("telephone = ?", telephone).First(&i)
	if i.Password == tools.MD5(password) {
		return nil
	}
	return errors.New("not equal")
}

func UserCheckAllow(target string, src string) bool {
	if target != src && !LinkExsit(target, src) {
		return false
	}
	return true
}
