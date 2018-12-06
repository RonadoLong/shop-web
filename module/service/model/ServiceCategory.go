package model

import "gopkg.in/mgo.v2/bson"

type ServiceCategory struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	EnName   string `json:"enName"`
	Imageurl string `json:"imageurl"`
	Type     string `json:"type"`
	Sort     int    `json:"sort"`
	ParentId int    `json:"parentId"`
	Status   int    `json:"status"`
}

type ServiceCLass struct {
	Id       bson.ObjectId `bson:"_id" json:"id"`
	Name     string        `json:"name"`
	EnName   string        `bson:"en_name" json:"enName"`
	ImgUrl   string        `bson:"img_url" json:"imgUrl"`
	Language string        `json:"language"`
	Settings []string      `json:"settings"`
	Status   int           `json:"status"`
}
