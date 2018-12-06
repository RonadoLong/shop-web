package home

import (
	"encoding/json"
	"os"
	"shop-web/common/commonUtils"
	"shop-web/common/cost"
	"shop-web/common/log"
	"shop-web/common/redis"
	model2 "shop-web/module/goods/model"
	"shop-web/module/home/dao"
	"shop-web/module/home/protocol"
	"strings"

	"github.com/gin-gonic/gin"
)

var Logger = log.NewLogger(os.Stdout)

var HomeService = homeService{}

type homeService struct {
}

func (*homeService) FindHomeNav(c *gin.Context) {

	var language = c.GetString("language")
	var err error

	homeHeaderReq := model.HomeHeaderResp{}
	key := model.HomeHeaderKey + language
	if exists, err := redis.RedisClient.Exists(key).Result(); err == nil {
		data := redis.RedisClient.Get(key).Val()
		if err == nil && exists != 0 {
			json.Unmarshal([]byte(data), &homeHeaderReq)
			commonUtils.CreateSuccess(c, homeHeaderReq)
			return
		}
	}

	if err != nil {
		Logger.Error(err)
		return
	}

	homeNavList, err := dao.FindHomeNav()
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	for idx := range homeNavList {
		if strings.Index(homeNavList[idx].ImgUrl, "http") == -1 {
			homeNavList[idx].ImgUrl = cost.Img_prefix + homeNavList[idx].ImgUrl
			if language == "EN" {
				homeNavList[idx].Title = homeNavList[idx].EnTitle
			}
		}
	}
	homeCarouselList, err := dao.FindHomeCarousel()
	var homeCarouselListReq = []model.HomeCarouselResp{}

	for _, homeCarousel := range homeCarouselList {
		homeCarouselReq := model.HomeCarouselResp{}
		if language == "EN" {
			homeCarouselReq.Title = homeCarousel.EnTitle
		} else {
			homeCarouselReq.Title = homeCarousel.Title
		}

		if strings.Index(homeCarousel.ImgUrl, "http") == -1 {
			homeCarouselReq.ImgUrl = cost.Img_prefix + homeCarousel.ImgUrl
		} else {
			homeCarouselReq.ImgUrl = homeCarousel.ImgUrl
		}
		homeCarouselReq.Id = homeCarousel.Id
		homeCarouselReq.Url = homeCarousel.Url
		homeCarouselListReq = append(homeCarouselListReq, homeCarouselReq)
	}

	homeHeaderReq.HomeCarouselList = homeCarouselListReq
	homeHeaderReq.HomeNavList = homeNavList

	rer, err := json.Marshal(homeHeaderReq)
	err = redis.RedisClient.Set(key, rer, 0).Err()

	commonUtils.CreateSuccess(c, homeHeaderReq)
}

func (*homeService) FindHomeList(c *gin.Context) {

	language := c.GetString("language")
	var err error

	key := model.HomeListKey + language
	var homeListResp = &model.HomeListResp{}

	// if _, err := redis.RedisClient.Exists(key).Result(); err == nil {
	// 	data, err := redis.RedisClient.Get(key).Result()
	// 	if err == nil {
	// 		json.Unmarshal([]byte(data), &homeListResp)
	// 		commonUtils.CreateSuccess(c, homeListResp)
	// 		return
	// 	}
	// }

	// if err != nil {
	// 	Logger.Error(err)
	// 	commonUtils.CreateError(c)
	// 	return
	// }

	//get news
	newsList, err := dao.FindHomeNews()
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}
	//get goods
	homeGoodsList, err := dao.FindHomeGoods()
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}
	goodsList := []model2.GoodsResp{}
	for _, homeGoods := range homeGoodsList {
		goods, _ := dao.FindGoodsByGoodsId(homeGoods.Id)
		goodsResp, err := commonUtils.TransferValue(goods, language)
		if err != nil {
			Logger.Error(err)
		}
		goodsList = append(goodsList, goodsResp)
	}

	homeListResp.HomeGoodsList = goodsList
	homeListResp.HomeNewsList = newsList

	rer, err := json.Marshal(homeListResp)
	err = redis.RedisClient.Set(key, rer, 0).Err()
	commonUtils.CreateSuccess(c, homeListResp)
}
