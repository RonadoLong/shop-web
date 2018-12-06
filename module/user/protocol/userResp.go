package protocol

import (
	"time"
)

type UserResp struct {
	UserId string  `json:"userId"  `
	Nickname string `json:"nickname" `
	RealName string `json:"realName" `
	Sex string `json:"sex" `
	Birth time.Time `json:"birth"`
	Avatar string `json:"avatar" `
	Hometown string `json:"hometown" `
	Remark string `json:"remark" `
	LoginTime time.Time `json:"loginTime" `
}
