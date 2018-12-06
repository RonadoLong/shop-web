package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type ServiceRoom struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Category     int           `json:"category"`
	UserId       string        `json:"userId"`
	Username     string        `json:"username"`
	ContactPhone string        `json:"contactPhone"`
	State        string        `json:"state"`
	City         string        `json:"city"`
	Address      string        `json:"address"`
	Area         string        `json:"area"`
	Code         string        `json:"code"`
	Location     GeoJson       `json:"location"`
	Pics         []string      `json:"pics"`
	Title        string        `json:"title"`
	Price        int           `json:"price"`
	DescStr      string        `json:"descStr"`
	Language     string        `json:"language"`
	Type         string        `json:"type"`
	Distance     float64       `bson:"distance" json:"distance"`
	Status       int           `json:"status"`
	UpdateAt     time.Time     `json:"updateAt"`
	CreateAt     time.Time     `json:"createAt"`
	LeaseType    int           `json:"leaseType"`
	RoomType     []string      `json:"roomType"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type ServiceJob struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Category     int           `json:"category"`
	UserId       string        `json:"userId"`
	Username     string        `json:"username"`
	ContactPhone string        `json:"contactPhone"`
	State        string        `json:"state"`
	City         string        `json:"city"`
	Address      string        `json:"address"`
	Area         string        `json:"area"`
	Code         string        `json:"code"`
	Location     GeoJson       `json:"location"`
	Pics         []string      `json:"pics"`
	Title        string        `json:"title"`
	Price        int           `json:"price"`
	DescStr      string        `json:"descStr"`
	Language     string        `json:"language"`
	Type         string        `json:"type"`
	Distance     float64       `bson:"distance" json:"distance"`
	Status       int           `json:"status"`
	UpdateAt     time.Time     `json:"updateAt"`
	CreateAt     time.Time     `json:"createAt"`
	Require      string        `json:"require"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type ServiceReq struct {
	//Longitude float64  `json:"longitude"`
	//Latitude  float64  `json:"latitude"`
	PageSize int      `json:"pageSize"`
	PageNum  int      `json:"pageNum"`
	State    []string `json:"state"`
	Category string   `json:"category"`
	RoomType []string `json:"roomType"`
}

type ServicePaymentSetting struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	EnName   string    `json:"en_name"`
	Time     int       `json:"time"`
	Price    int       `json:"price"`
	Status   int       `json:"status"`
	UpdateAt time.Time `json:"updateAt"`
	CreateAt time.Time `json:"createAt"`
}

type ServicePayFlow struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ServiceId string        `json:"serviceId"`
	ExTime    time.Time     `json:"exTime"`
	Price     int           `json:"price"`
	Status    int           `json:"status"`
	UpdateAt  time.Time     `json:"updateAt"`
	CreateAt  time.Time     `json:"createAt"`
}

type ServicePayFlowProtocol struct {
	ServiceId string
	ExTime    int
	Price     int
}
