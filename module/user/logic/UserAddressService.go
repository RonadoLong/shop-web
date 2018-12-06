package logic

import (
	"github.com/gin-gonic/gin"
	"shop-web/module/user/model"
	"shop-web/common/dbutil"
	"shop-web/common/commonUtils"
)

func (userService * userService) SaveUserAddress(c *gin.Context) {
	if userId := c.GetString("userId"); userId != "" {

		var userAddress  model.UserAddress
		if err := c.ShouldBindJSON(&userAddress); err != nil{
			logger.Error(err)
			commonUtils.CreateErrorParams(c)
			return
		}

		userAddress.UserId = userId
		if 	err := dbutil.DB.Table("user_address").Save(&userAddress).Error; err != nil {
			logger.Error(err)
		}else {
			commonUtils.CreateSuccess(c, nil)
			return
		}
	}

	commonUtils.CreateErrorParams(c)
}

func (userService * userService) FindAddressByUserId(c *gin.Context) {
	if userId := c.GetString("userId"); userId != "" {
		 var  userAddressList []model.UserAddress
		if 	err := dbutil.DB.Table("user_address").
			Where("user_id = ? and status = 1 ", userId).
			Order("`create_at` desc").Find(&userAddressList).Error; err != nil {
			logger.Error(err)
		}else {

			if len(userAddressList) == 0 {
				commonUtils.CreateNotContent(c)
				return
			}

			 commonUtils.CreateSuccess(c, userAddressList)
			return
		}
		//userAddressList := make([]model.UserAddress, 10)
	}

	commonUtils.CreateErrorParams(c)
}