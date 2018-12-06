package model

import (
	"time"
	)

type ServiceOrder struct {
	OrderId string `json:"orderId"`
	UserId string `json:"user_id"`
	ServiceId string `json:"serviceId"`
	AutoPay int `json:"autoPay"`
	IsRenew int `json:"isRenew"`
	PayPrice int `json:"payPrice"`
	PayStatus int `json:"payStatus"`
	StartTime time.Time `json:"startTime"`
	ExpireTime time.Time `json:"expireTime"`
	Status int `json:"status"`
	UpdateAt time.Time `json:"updateAt"`
	CreateAt time.Time `json:"createAt"`
}
