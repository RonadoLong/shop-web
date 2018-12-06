package controller

import (
	"github.com/gin-gonic/gin"
	"shop-web/module/news/service"
)

var VideoRest = videoRest{}

type videoRest struct {
}

func (*videoRest)FindVideoList(c *gin.Context)  {
	service.VideoService.FindVideoList(c)
}

func (*videoRest)FindVideoAllList(c *gin.Context)  {
	service.VideoService.FindVideoAllList(c)
}

func (*videoRest)FindVideoDetail(c *gin.Context){
	service.VideoService.FindVideoDetail(c)
}

func (*videoRest) FindCategoryList(c *gin.Context){
	service.VideoService.FindCategoryList(c)
}