package logic

import (
	"github.com/gin-gonic/gin"
	"shop-web/module/user/protocol"
	"shop-web/common/redis"
	"shop-web/module/user/utils"
	"shop-web/module/user/dao"
	"os"
	"shop-web/common/log"
	"shop-web/module/user/model"
	"time"
	"strconv"
	"sync"
	"github.com/HaroldHoo/id_generator"
	"shop-web/common/dbutil"
	"shop-web/common/commonUtils"
)

var logger = log.NewLogger(os.Stdout)

var UserService = &userService{
	mutex: &sync.Mutex{},
}

type userService struct {
	mutex *sync.Mutex
}

func (userService *userService)PhoneLogin(c *gin.Context, login *protocol.LoginReq) (*model.User, bool) {

	key := utils.UserCodeKey + login.Phone

	if  _, err := redis.RedisClient.Exists(key).Result(); err == nil {
		code, err := redis.RedisClient.Get(key).Result()
		redis.RedisClient.Del(key)
		if err == nil {
			if login.VerifyCode == string(code) {
				//查询用户是否存在
				userAuth := dao.HasSameUserByUnionId("phone", login.Phone)
				if userAuth != nil && userAuth.UserId != ""{
					user, err := dao.FindUserByUserId(userAuth.UserId)
					if err == nil {
						user.IsBindPhone = true
						logger.Infof("user is exit %s", user.UserId)
						return user, true
					}
					logger.Error(err)
				}
				//create new user
				userService.mutex.Lock()
				defer userService.mutex.Unlock()

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
						Identify: login.Phone,
						UpdateTime: time.Now(),
						CreateTime: time.Now(),
						Status: "1",
					}
					user := utils.CreateIphoneLoginUser()
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
			}
		}
	}

	return nil, false
}

func (userService *userService) GetUserIntegralFlow(context *gin.Context) {
	if userId := context.GetString("userId"); userId != "" {
		var  userIntegralFlow []model.UserIntegralFlow
		if 	err := dbutil.DB.Table("user_integral_flow").Where("user_id = ? and r_user_id is not null and r_user_id != '' ", userId).
			Order("`create_at` desc").Find(&userIntegralFlow).Error; err != nil {
			logger.Error(err)
		}else {

			if len(userIntegralFlow) == 0 {
				commonUtils.CreateNotContent(context)
				return
			}
			commonUtils.CreateSuccessByList(context, len(userIntegralFlow),userIntegralFlow)
			return
		}
	}

}

