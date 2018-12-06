package model

import (
	"time"
)

type UserAuth struct {
	Id string `json:"id" gorm:"primary_key"`
	UserId string `json:"userId" `
	IdentifyType string `json:"identifyType" `
	Identify string `json:"identify" `
	Credential string `json:"credential" `
	TradeId string`json:"tradeId" `
	Status string `json:"status" `
	CreateTime time.Time `json:"createTime" `
	UpdateTime time.Time `json:"updateTime" `
}
