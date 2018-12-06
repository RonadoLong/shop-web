package service

import (
	"sync"
	"github.com/gin-gonic/gin"
	"os"
	"shop-web/common/log"
	"shop-web/module/news/dao"
	"shop-web/common/commonUtils"
	"strconv"
	"shop-web/module/news/newsModel"
	"shop-web/common/dbutil"
	"shop-web/common/redis"
)

var Logger = log.NewLogger(os.Stdout)

var NewsService = newsService{
	mutex: &sync.Mutex{},
}
type newsService struct {
	mutex *sync.Mutex
}


func (*newsService)FindNewsList(c *gin.Context)  {

	offset,pageSize := commonUtils.GetOffset(c)
	category := "[\"" + c.Param("category") + "\"]"
	Logger.Info(category, pageSize)
	total,err := dao.FindNewsCount(category)
	if err != nil{
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}
	if offset >= total{
		commonUtils.CreateNotContent(c)
		return
	}
	newsList,err := dao.FindNewsListByOffset(offset, pageSize, category)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}
	commonUtils.CreateSuccessByList(c,total,newsList)
}

func (service *newsService) FindNewsAllList(c *gin.Context) {
	offset,pageSize := commonUtils.GetOffset(c)
	total,err := dao.FindNewsALLCount()
	if err != nil{
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}
	if offset >= total{
		commonUtils.CreateNotContent(c)
		return
	}
	newsList,err := dao.FindNewsAllListByOffset(offset, pageSize)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}
	commonUtils.CreateSuccessByList(c,total,newsList)
}


func (*newsService)FindNewsDetail(c *gin.Context) {

	var err error
	newsId,_ := strconv.ParseInt(c.Param("newsId"),10, 64)
	if newsId >= 1 {
		news,err := dao.FindNewsDetail(newsId)
		if err == nil {
			commonUtils.CreateSuccess(c, news)
			return
		}
	}

	Logger.Error(err)
	commonUtils.CreateError(c)
}


func (*newsService)FindNewSNavList(c *gin.Context) {

	var err error
	defer Logger.Error(err)

	var cateList []newsModel.NewsCategory
	err = dbutil.DB.Table("news_category").Order("`sort` asc").Where("`status` = 1 ").Find(&cateList).Error
	if err != nil{
		commonUtils.CreateError(c)
		return
	}
	commonUtils.CreateSuccess(c,cateList)
}


func (*newsService)AddLikeById(c *gin.Context) {

	var err error
	defer Logger.Error(err)

	newsLikeKey := redis.NEWS_LIKE_KEY + c.Param("newsId")
	ret, err := redis.RedisClient.Exists(newsLikeKey).Result()
	if ret > int64(0){
		redis.RedisClient.IncrBy(c.Param("newsId"), int64(1))
	}else {
		redis.RedisClient.Set(newsLikeKey, 1,0)
	}

	if err != nil{
		commonUtils.CreateError(c)
		return
	}

	commonUtils.CreateSuccess(c,nil)

}





