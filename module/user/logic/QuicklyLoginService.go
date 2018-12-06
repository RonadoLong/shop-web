package logic

import (
	"github.com/gin-gonic/gin"
	"shop-web/module/user/protocol"
	"shop-web/module/user/model"
	"shop-web/module/user/dao"
	"strconv"
	"time"
	"shop-web/module/user/utils"
	"github.com/HaroldHoo/id_generator"
)

func (userService *userService)FaceBookLogin(c *gin.Context, login *protocol.LoginReq) (*model.User, bool) {

	userService.mutex.Lock()
	defer userService.mutex.Unlock()

	// user is exits ?
	userAuth := dao.HasSameUserByUnionId(login.Type, login.UnionId)
	if userAuth != nil && userAuth.UserId != ""{
		user,err := dao.FindUserByUserId(userAuth.UserId)
		if err == nil {
			user.IsBindPhone = true
			logger.Infof("user is exit %s", user.UserId)
			return user, true
		}
		logger.Error(err)
	}

	// Data ID, like Mysql table ID, 10bit （设置10bit的数据ID）
	var dataID = uint64(256)
	id, gErr := id_generator.NextId(dataID)
	if gErr == nil {
		userIdStr := strconv.FormatUint(id, 10)
		recommendCode := strconv.FormatUint(id, 36)
		userAuth := &model.UserAuth{
			Id: userIdStr,
			UserId: userIdStr,
			IdentifyType: login.Type,
			Identify: login.UnionId,
			UpdateTime: time.Now(),
			CreateTime: time.Now(),
			Status: "1",
		}
		//create new user
		user := utils.CreateQuicklyLoginUser(login)
		user.UserId = userIdStr
		user.RecommendCode = recommendCode

		if  login.RecommendCode != ""{
			ruser := dao.FindUserByRecommendCode(login.RecommendCode)
			if ruser != nil{
				recommend := &model.UserRecommend{
					RecommendCode: ruser.RecommendCode,
					RecommendUserId: ruser.UserId,
					UserId: userIdStr,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}
				err := dao.CreateUser(userAuth, user, recommend)
				if err == nil {
					return user, true
				}
			}
		}else {
			err := dao.CreateUser(userAuth, user, nil)
			logger.Error(err)
			if err == nil {
				return user, true
			}
		}
	}
	return nil, false
}