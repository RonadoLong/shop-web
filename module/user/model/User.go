package model

import (
	"shop-web/common/commonModel"
	"time"
)

type User struct {
	UserId string  `json:"userId" gorm:"primary_key;size:20" `
	Nickname string `json:"nickname" gorm:"size:20"`
	RealName string `json:"realName" gorm:"size:20"`
	Sex string `json:"sex" gorm:"type:char;size:1"`
	Birth time.Time `json:"birth"`
	Avatar string `json:"avatar" `
	Hometown string `json:"hometown" `
	Remark string `json:"remark" `
	LoginTime time.Time `json:"loginTime" `
	IsRecommend string `json:"isRecommend" gorm:"NOT NULL;DEFAULT:0"`
	RecommendCode string `json:"recommendCode" gorm:"NOT NULL;unique_index"`
	Duration int64 `json:"duration" gorm:"DEFAULT:0"`
	Integral int `json:"integral"  gorm:"DEFAULT:0"`
	Commission int `json:"commission" gorm:"DEFAULT:0.00"`
	Status string `json:"status"  gorm:"type:char;size:1"`
	commonModel.Model
	IsBindWx bool `json:"isBindWx"  gorm:"-"`
	IsBindPhone bool `json:"isBindPhone"  gorm:"-"`
}
