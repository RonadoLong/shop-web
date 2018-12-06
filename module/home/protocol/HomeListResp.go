package model

import (
	"shop-web/module/news/newsModel"
	"shop-web/module/goods/model"
)

type HomeListResp struct {
	HomeNewsList [] newsModel.News
	HomeGoodsList []model.GoodsResp
}