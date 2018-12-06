package controller

import (
	"github.com/gin-gonic/gin"
	"shop-web/module/user/protocol"
	"os"
	"shop-web/common/log"
	"shop-web/module/user/model"
	"shop-web/module/user/logic"
	"shop-web/common/middleware"
)

var logger = log.NewLogger(os.Stdout)

var UserRest = &userController{}

type userController struct {
}

func (*userController) Login(login interface{}, c *gin.Context) (*model.User, bool) {
	if login != nil {
		var loginReq = login.(protocol.LoginReq)
		if loginReq.Type == "wechat" || loginReq.Type == "facebook" {
			return logic.UserService.FaceBookLogin(c, &loginReq)
		} else if loginReq.Type == "phone" {
			return logic.UserService.PhoneLogin(c, &loginReq)
		}
	}
	return nil, false
}

func (*userController) SendCode(c *gin.Context) {
	middleware.MessageService.SendCode(c)
}

func (*userController) SaveAgreement(c *gin.Context) {
	logic.UserService.SaveUserAgreement(c)
}

func (*userController) GetUserAgreement(c *gin.Context) {
	logic.UserService.GetUserAgreement(c)
}

func (*userController) GetUserAddress(c *gin.Context) {
	logic.UserService.FindAddressByUserId(c)
}

func (*userController) SaveUserAddress(c *gin.Context) {
	logic.UserService.SaveUserAddress(c)
}
