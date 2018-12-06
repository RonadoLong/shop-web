package model

import "time"

type UserRecommend struct {
	Id              int64     `json:"id" gorm:"primary_key;AUTO_INCREMENT;UNIQUE_INDEX"`
	RecommendUserId string    `json:"recommendUserId" `
	UserId          string    `json:"userId"`
	RecommendCode   string    `json:"recommendCode"`
	CreateTime      time.Time `json:"createTime"`
	UpdateTime      time.Time `json:"updateTime"`
	Status          string    `json:"status" gorm:"DEFAULT"`
}

type UserIntegralFlow struct {
	Id         int64     `json:"id"`
	UserId     string    `json:"userId"`
	RUserId    string    `json:"rUserId"`
	OrderId    string    `json:"orderId"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
	Integral   int       `json:"integral"`
}
