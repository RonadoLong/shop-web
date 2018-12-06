package controller

import (
	"github.com/gin-gonic/gin"
	"shop-web/common/commonUtils"
	"shop-web/module/user/logic"
	"shop-web/module/user/protocol"
	"github.com/gin-gonic/gin/binding"
	"time"
)

func (*userController) GetUserInfo(c *gin.Context)  {
	if userId, isExits := c.Get("userId"); isExits {
		logger.Info(userId.(string))
		userInfo, err := logic.UserService.GetUserInfo(userId.(string))
		if err == nil {
			commonUtils.CreateSuccess(c, userInfo)
			return
		}
		logger.Error(err)
	}
	commonUtils.CreateErrorParams(c)
}

func (*userController) UpdateUserInfo(c *gin.Context)  {
	var userReq = protocol.UserResp{}
	err := c.ShouldBindWith(&userReq, binding.JSON)
	if err == nil {

		if userId, exits := c.Get("userId"); exits == true {
			userReq.UserId = userId.(string)
			userReq.LoginTime = time.Now()
			err := logic.UserService.UpdateUserInfo(userReq)
			if err == nil {
				commonUtils.CreateSuccess(c, nil)
				return
			}
		}
	}
	logger.Error(err)
	commonUtils.CreateErrorParams(c)
}

func (*userController) GetUserIntegralFlow(c *gin.Context)  {
	logic.UserService.GetUserIntegralFlow(c)
}
