package logic

import (
	"shop-web/common/commonUtils"
	"shop-web/common/dbutil"
	"shop-web/module/user/dao"
	"shop-web/module/user/model"
	"shop-web/module/user/protocol"

	"github.com/gin-gonic/gin"
)

func (userService *userService) GetUserInfo(userId string) (interface{}, error) {
	return dao.FindUserByUserId(userId)
}

func (userService *userService) UpdateUserInfo(resp protocol.UserResp) error {
	return dao.UpdateUserInfo(resp)
}

func (userService *userService) GetUserAgreement(c *gin.Context) {
	if userId, exits := c.Get("userId"); exits == true && userId != "" {
		var count int
		if err := dbutil.DB.Table("user_agreement").Where(" userId = ? ", userId).Count(&count).Error; err == nil {
			commonUtils.CreateSuccess(c, count)
			return
		}
	}
	commonUtils.CreateErrorParams(c)
}

func (userService *userService) SaveUserAgreement(c *gin.Context) {
	userService.mutex.Lock()
	defer userService.mutex.Unlock()

	if userId := c.GetString("userId"); userId != "" {
		userAgreement := model.UserAgreement{}
		userAgreement.UserId = userId
		userAgreement.Status = 1

		if err := dbutil.DB.Table("user_agreement").Save(&userAgreement).Error; err != nil {
			logger.Error(err)
		} else {
			commonUtils.CreateSuccess(c, nil)
			return
		}
	}
	commonUtils.CreateErrorParams(c)
}



