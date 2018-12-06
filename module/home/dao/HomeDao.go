package dao

import (
	"shop-web/module/home/model"
	"shop-web/common/dbutil"
	GoodsModel "shop-web/module/goods/model"
	"shop-web/module/news/newsModel"
)

func FindHomeNav() ([]model.HomeNav, error) {
	var homeNavList []model.HomeNav
	err := dbutil.DB.Table("home_nav").Where("`status` = 1").Order("`sort` desc").Find(&homeNavList).Error
	if err != nil{
		return homeNavList, err
	}
	return homeNavList, nil
}

func FindHomeCarousel() ([]model.HomeCarousel,error){
	var homeCarousel []model.HomeCarousel
	err := dbutil.DB.Table("home_carousel").Where("`status` = 1").Order("`sort` desc").Find(&homeCarousel).Error
	if err != nil{
		return homeCarousel, err
	}
	return homeCarousel, nil
}

func FindHomeGoods() ([]model.HomeGoods, error) {
	var homeGoodsList []model.HomeGoods
	err := dbutil.DB.Table("home_goods").Where("`status` = 1").Find(&homeGoodsList).Error
	if err != nil {
		return nil, err
	}
	return homeGoodsList, nil
}

func FindGoodsByGoodsId(id int64) ( GoodsModel.Goods, error)  {
	var goods GoodsModel.Goods
	err := dbutil.DB.Table("goods").Where("`goods_id` = ? ", id).Find(&goods).Error
	if err != nil {
		return goods,err
	}
	return goods, nil
}

func FindGoodsDetailByGoodsId(id int64) ( GoodsModel.GoodsDetail, error)  {
	var goodsDetail GoodsModel.GoodsDetail
	err := dbutil.DB.Table("goods_detail").Where("`goods_id` = ? ", id).Find(&goodsDetail).Error
	if err != nil {
		return goodsDetail,err
	}
	return goodsDetail, nil
}

func FindGoodsSkuByGoodsId(id int64) ([]GoodsModel.GoodsSku, error)  {
	var goodsSku []GoodsModel.GoodsSku
	err := dbutil.DB.Table("goods_sku").Where("`goods_id` = ? ", id).Find(&goodsSku).Error
	if err != nil {
		return goodsSku,err
	}
	return goodsSku, nil
}


func FindGoodsSkuPropertyBySkuId(id int64) ([]GoodsModel.GoodsSkuProperty, error)  {
	var goodsSkuPropertys  []GoodsModel.GoodsSkuProperty
	err := dbutil.DB.Table("goods_sku_property").Where("`sku_id` = ? ", id).Find(&goodsSkuPropertys).Error
	if err != nil {
		return goodsSkuPropertys,err
	}
	return goodsSkuPropertys, nil
}

func FindHomeNews() ([]newsModel.News, error)  {
	var News []newsModel.News
	err := dbutil.DB.Table("home_news").Find(&News).Error
	if err != nil {
		return News,err
	}
	return News, nil
}

//
//func FindNewsById(id int64) (newsModel.News, error)  {
//	var News newsModel.News
//	err := dbutil.DB.Table("news").Where(" `id` = ? and `id` IS NOT NULL", id).Find(&News).Error
//	if err != nil {
//		return News,err
//	}
//	return News, nil
//}
