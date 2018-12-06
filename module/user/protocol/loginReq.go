package protocol

import "time"

type LoginReq struct {
	Username string `form:"username" json:"username" `
	Password string `form:"password" json:"password" `
	UnionId string `form:"unionId" json:"unionId" `
	Type string `form:"type" json:"type" ` /** wechat  2 facebook 3 phone*/
	Sex string `form:"sex" json:"sex" `
	Nickname string `form:"nickname" json:"nickname" `
	Birthday time.Time `form:"birthday" json:"birthday" `
	VerifyCode string `form:"verifyCode" json:"verifyCode" `
	Phone string `form:"phone" json:"phone" `
	RecommendCode string `form:"recommendCode" json:"recommendCode" `
	Avatar string `json:"avatar"`
	Hometown string `json:"hometown" `
}
