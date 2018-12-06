package dao

import (
	"shop-web/module/user/model"
	"shop-web/common/dbutil"
	"shop-web/common/log"
	"os"
	"shop-web/module/user/protocol"
)

var logger = log.NewLogger(os.Stdout)

func HasSameUserByUnionId(IdentifyType string, Identify string)  (*model.UserAuth){
	userAuth := model.UserAuth{}
	rows, err := dbutil.DB.Raw("select * from user_auth where identify_type = ? and identify = ? ", IdentifyType, Identify).Rows()
	if err == nil {
		for rows.Next() {
			dbutil.DB.ScanRows(rows, &userAuth)
		}
		return &userAuth
	}
	logger.Error(err)
	return nil
}

func FindUserByUserId(userId string) (*model.User,error){
	user := model.User{}
	rows, err := dbutil.DB.Raw("select * from user where user_id = ?  ", userId).Rows()
	if err == nil {
		for rows.Next() {
			dbutil.DB.ScanRows(rows, &user)
		}
		return &user, nil
	}
	return nil,err
}

func FindUserByRecommendCode(code string) *model.User {
	user := &model.User{}
	err := dbutil.DB.Table("user").Where("`recommend_code` = ?", code).Scan(user).Error
	if err != nil {
		logger.Errorf("get recommend_code user error %s ", err)
		return nil
	}
	return user
}

func CreateUser(userAuth *model.UserAuth, user *model.User, userRecommend *model.UserRecommend) error {

	tx := dbutil.DB.Begin()
	if err := tx.Table("user_auth").Create(userAuth).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Table("user").Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if userRecommend != nil{
		if err := tx.Table("user_recommend").Create(userRecommend).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	tx.Commit()
	return nil
}

func UpdateUserInfo(resp protocol.UserResp) error {
	tx := dbutil.DB.Begin()
	err := tx.Table("user").Update(&resp).Error
	if err != nil{
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetRecommendUserIdByUserId(userId string) (ret model.UserRecommend, err error){
	err = dbutil.DB.Exec("select * from user_recommend where user_id = ? ", userId).Find(&ret).Error
	if err != nil {
		return ret, err
	}
	return ret, err
}


func UpdateIntegralByUserId(userId string, totalIntegral int) (err error){
	err = dbutil.DB.Exec("update user set integral = integral + ?, update_time = now() where user_id = ?", totalIntegral, userId).Error
	return err
}

func SaveUserIntegralFlow(flow model.UserIntegralFlow) (err error){
	err = dbutil.DB.Table("user_integral_flow").Save(&flow).Error
	return err
}