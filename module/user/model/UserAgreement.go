package model

type UserAgreement struct {
	Id string `json:"id"`
	UserId string `json:"userId"`
	Status int `json:"status"`
}