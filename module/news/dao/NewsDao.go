package dao

import (
	"shop-web/common/dbutil"
	"shop-web/module/news/newsModel"
)

func FindNewsCount(category string) (int, error) {
	count := 0
	if err := dbutil.DB.Table("news").Where("category = ? and status = 1", category).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func FindNewsALLCount() (int, error) {
	count := 0
	if err := dbutil.DB.Table("news").Where("status = 1").Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func FindNewsAllListByOffset(pageNum int, pageSize int) ([]newsModel.News, error) {
	var newsList []newsModel.News

	if err := dbutil.DB.Table("news").Where("status = 1 ").Order("`create_time` desc").Offset(pageNum).Limit(pageSize).Find(&newsList).Error; err != nil {
		return newsList, err
	}
	return newsList, nil
}

func FindNewsListByOffset(pageNum int, pageSize int, category string) ([]newsModel.News, error) {
	var newsList []newsModel.News

	if err := dbutil.DB.Table("news").Where(" category = ? and status = 1 ", category).Order("`create_time` desc").Offset(pageNum).Limit(pageSize).Find(&newsList).Error; err != nil {
		return newsList, err
	}
	return newsList, nil
}

func FindNewsDetail(id int64) (newsModel.News, error) {
	var news newsModel.News
	if err := dbutil.DB.Table("news").Where(" `id` = ? and `status` = 1 ", id).Find(&news).Error; err != nil {
		return news, err
	}
	return news, nil
}
