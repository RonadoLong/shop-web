package controller

import (
	"github.com/gin-gonic/gin"
	"shop-web/common/log"
	"os"
	"shop-web/module/news/service"
)


var Logger = log.NewLogger(os.Stdout)
var NewsRest = newsRest{}

type newsRest struct {

}

func (newsRest *newsRest)FindNewsList(c *gin.Context){
	service.NewsService.FindNewsList(c)
}
func (newsRest *newsRest)FindNewsAllList(c *gin.Context){
	service.NewsService.FindNewsAllList(c)
}

func (newsRest *newsRest)FindNewsDetail(c *gin.Context){
	service.NewsService.FindNewsDetail(c)
}

func (newsRest *newsRest)FindNewSNavList(c *gin.Context){
	service.NewsService.FindNewSNavList(c)
}

func (newsRest *newsRest)AddLikeById(c *gin.Context){
	service.NewsService.AddLikeById(c)
}