package dao

import (
	"shop-web/common/dbutil"
	"shop-web/module/news/newsModel"
)


func FindVideoCount(category string) (int,error){
	count := 0
	if err := dbutil.DB.Table("video").Where("category = ? and status = 1", category).Count(&count).Error; err != nil{
		return count, err
	}
	return count, nil
}

func FindAllVideoCount() (int, error) {
	count := 0
	if err := dbutil.DB.Table("video").Where("status = 1").Count(&count).Error; err != nil{
		return count, err
	}
	return count, nil
}

func FindVideosAllListByOffset(pageNum int, pageSize int) ([]newsModel.Video,error){
	var videoList []newsModel.Video
	if err := dbutil.DB.Table("video").Where("status = 1").Order("`create_time` desc").Offset(pageNum).Limit(pageSize).Find(&videoList).Error; err != nil{
		return videoList, err
	}
	return videoList, nil
}


func FindVideosListByOffset(pageNum int, pageSize int, category string) ([]newsModel.Video,error){
	var videoList []newsModel.Video
	if err := dbutil.DB.Table("video").Where("category = ? and status = 1", category).Order("`create_time` desc").Offset(pageNum).Limit(pageSize).Find(&videoList).Error; err != nil{
		return videoList, err
	}
	return videoList, nil
}

func FindVideoDetail(id int64)(newsModel.Video,error){
	var video newsModel.Video
	if err := dbutil.DB.Table("video").Where(" `id` = ? and `status` = 1 ", id).Find(&video).Error; err != nil{
		return video, err
	}
	return video,nil
}