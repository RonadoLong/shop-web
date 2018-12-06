package model

import (
	"gopkg.in/mgo.v2/bson"
)

type ServiceResp struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ClassId  int           `json:"classId"`
	UserId   string        `json:"userId"`
	Name     string        `json:"name"`
	Tel      string        `json:"tel"`
	Cell     string        `json:"cell"`
	Fax      string        `json:"fax"`
	Wechat   string        `json:"wechat"`
	Email    string        `json:"email"`
	State    []string      `json:"state"`
	City     string        `json:"city"`
	Area     string        `json:"area"`
	Location []float64     `json:"location"`
	RoomType []string      `json:"roomType"`
	Banner   []string      `json:"banner"`
	Title    string        `json:"title"`
	Price    int           `json:"price"`
	Note     string        `json:"note"`
}

type ServiceJobResp struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Category     int           `json:"category"`
	UserId       string        `json:"userId"`
	Require      string        `json:"require"`
	Username     string        `json:"username"`
	ContactPhone string        `json:"contactPhone"`
	State        string        `json:"state"`
	City         string        `json:"city"`
	Address      string        `json:"address"`
	Area         string        `json:"area"`
	Code         string        `json:"code"`
	Location     []float64     `json:"location"`
	Pics         []string      `json:"pics"`
	Title        string        `json:"title"`
	Price        int           `json:"price"`
	Type         string        `json:"type"`
	DescStr      string        `json:"descStr"`
	Language     string        `json:"language"`
	Status       int           `json:"status"`
}
