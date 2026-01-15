package database

import (
	"errors"

	"github.com/dingdinglz/ai-swindle-detecter-backend/tools"
)

type UserInfo struct {
	UserID    string
	Telephone string
	Nickname  string
	Avatar    string
	Email     string
}

func UserNew(telephone string, password string) (*UserInfo, error) {
	if UserExist(telephone) {
		return nil, errors.New("user has exists")
	}
	hashedPassword, err := tools.HashPassword(password)
	if err != nil {
		return nil, err
	}
	userID := "user_" + telephone
	var i UserTable
	i.Telephone = telephone
	i.Password = hashedPassword
	i.UserID = userID
	i.Nickname = ""
	i.Avatar = ""
	i.Email = ""
	MainDB.Create(&i)
	return &UserInfo{
		UserID:    userID,
		Telephone: telephone,
		Nickname:  "",
		Avatar:    "",
		Email:     "",
	}, nil
}

func UserExist(telephone string) bool {
	var i int64
	MainDB.Model(&UserTable{}).Where("telephone = ?", telephone).Count(&i)
	return i != 0
}

func UserLogin(telephone string, password string) (*UserInfo, error) {
	var i UserTable
	result := MainDB.Model(&UserTable{}).Where("telephone = ?", telephone).First(&i)
	if result.Error != nil {
		return nil, errors.New("not exist")
	}
	if tools.CheckPassword(password, i.Password) {
		return &UserInfo{
			UserID:    i.UserID,
			Telephone: i.Telephone,
			Nickname:  i.Nickname,
			Avatar:    i.Avatar,
			Email:     i.Email,
		}, nil
	}
	return nil, errors.New("not equal")
}

func UserCheckAllow(target string, src string) bool {
	if target != src && !LinkExsit(target, src) {
		return false
	}
	return true
}

func GetUserByTelephone(telephone string) (*UserInfo, error) {
	var i UserTable
	result := MainDB.Model(&UserTable{}).Where("telephone = ?", telephone).First(&i)
	if result.Error != nil {
		return nil, errors.New("not exist")
	}
	return &UserInfo{
		UserID:    i.UserID,
		Telephone: i.Telephone,
		Nickname:  i.Nickname,
		Avatar:    i.Avatar,
		Email:     i.Email,
	}, nil
}

