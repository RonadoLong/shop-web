package service

import (
	"sync"
	"github.com/gin-gonic/gin"
	"shop-web/module/news/dao"
	"shop-web/common/commonUtils"
	"strconv"
	"shop-web/common/dbutil"
	"shop-web/module/news/newsModel"
	"encoding/json"
)



var VideoService = videoService{
	mutex: &sync.Mutex{},
}

type videoService struct {
	mutex *sync.Mutex
}


func (*videoService)FindVideoList(c *gin.Context)  {

	category := c.Param("category")
	offset,pageSize := commonUtils.GetOffset(c)
	Logger.Info(category, pageSize)

	if category == ""{
		 commonUtils.CreateErrorParams(c)
		 return
	}

	total,err := dao.FindVideoCount(category)
	if err != nil{
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	if offset >= total{
		commonUtils.CreateNotContent(c)
		return
	}

	videoList,err := dao.FindVideosListByOffset(offset, pageSize, category)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	videoRespList := []newsModel.VideoResp{}
	for _,video := range videoList {
		videoResp := newsModel.VideoResp{}
		videoResp.Id = video.Id
		videoResp.Title = video.Title
		videoResp.ThumbUrl = video.ThumbUrl
		videoResp.Duration = video.Duration
		videoResp.ReadCount = video.ReadCount
		videoResp.CommentCount = video.CommentCount
		videoResp.LikeCount = video.LikeCount
		videoResp.Category = video.Category
		videoResp.Content = video.Content
		videoResp.Tags = video.Tags
		videoResp.VideoDesc = video.VideoDesc

		var pusherInfo newsModel.PusherInfo
		err := json.Unmarshal([]byte(video.PusherInfo), &pusherInfo)
		if err == nil {
			videoResp.Author = pusherInfo.Name
			videoResp.Avatar = pusherInfo.Avatar
		}
		videoRespList = append(videoRespList, videoResp)
	}
	commonUtils.CreateSuccessByList(c, total, videoRespList)
}

func (service *videoService) FindVideoAllList(c *gin.Context) {
	offset,pageSize := commonUtils.GetOffset(c)
	total,err := dao.FindAllVideoCount()
	if err != nil{
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	if offset >= total{
		commonUtils.CreateNotContent(c)
		return
	}

	videoList,err := dao.FindVideosAllListByOffset(offset, pageSize)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	videoRespList := []newsModel.VideoResp{}
	for _,video := range videoList {
		videoResp := newsModel.VideoResp{}
		videoResp.Id = video.Id
		videoResp.Title = video.Title
		videoResp.ThumbUrl = video.ThumbUrl
		videoResp.Duration = video.Duration
		videoResp.ReadCount = video.ReadCount
		videoResp.CommentCount = video.CommentCount
		videoResp.LikeCount = video.LikeCount
		videoResp.Category = video.Category
		videoResp.Content = video.Content
		videoResp.Tags = video.Tags
		videoResp.VideoDesc = video.VideoDesc

		var pusherInfo newsModel.PusherInfo
		err := json.Unmarshal([]byte(video.PusherInfo), &pusherInfo)
		if err == nil {
			videoResp.Author = pusherInfo.Name
			videoResp.Avatar = pusherInfo.Avatar
		}
		videoRespList = append(videoRespList, videoResp)
	}
	commonUtils.CreateSuccessByList(c, total, videoRespList)
}

func (*videoService)FindVideoDetail(c *gin.Context) {

	var err error
	defer Logger.Error(err)

	videoId,_ := strconv.ParseInt(c.Param("videoId"),10, 64)
	if videoId >= 1 {
		news,err := dao.FindVideoDetail(videoId)
		if err == nil {
			commonUtils.CreateSuccess(c, news)
			return
		}
	}

	commonUtils.CreateError(c)
}

func (*videoService)FindCategoryList(c *gin.Context)  {
	var err error
	defer Logger.Error(err)

	var cateList []newsModel.VideoCategory
	err = dbutil.DB.Table("video_category").Order("`sort` asc").Where("`status` = 1 ").Find(&cateList).Error
	if err != nil{
		commonUtils.CreateError(c)
		return
	}
	commonUtils.CreateSuccess(c,cateList)
}


