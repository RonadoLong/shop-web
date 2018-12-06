package utils

import (
	"shop-web/module/user/model"
	"time"
	"shop-web/module/user/protocol"
)

func CreateIphoneLoginUser() *model.User{
	user := model.User{}
	user.Nickname = "user_tohnet"
	user.Sex = "1"
	user.Avatar = "http://images.samecity.com.cn/user/default/default_image.png"
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	user.LoginTime = time.Now()
	user.Birth =  time.Now()
	user.Remark = "welcome to tohnet"
	user.Status = "1"
	return &user
}

func CreateQuicklyLoginUser(req *protocol.LoginReq) *model.User {
	user := model.User{}
	user.Sex = "0"
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	user.LoginTime = time.Now()
	user.Remark = "welcome to tohnet"
	user.Birth =  time.Now()
	user.Status = "1"
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}else {
		user.Avatar = "http://images.samecity.com.cn/user/default/default_image.png"
	}

	if  req.Nickname != ""{
		user.Nickname = req.Nickname
	}else {
		user.Nickname = "user_tohnet"
	}

	if req.Sex != "" {
		user.Sex = req.Sex
	}

	if  req.Hometown != ""{
		user.Hometown = req.Hometown
	}

	return &user
}