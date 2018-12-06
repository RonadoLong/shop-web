package model

import "shop-web/module/home/model"

type HomeHeaderResp struct {
	HomeNavList []model.HomeNav `json:"homeNavList"`
	HomeCarouselList []HomeCarouselResp `json:"homeCarouselList"`
}

type HomeCarouselResp struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	ImgUrl string `json:"imgUrl"`
	Url string `json:"url"`
}