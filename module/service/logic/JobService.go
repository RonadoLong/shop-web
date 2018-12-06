package logic

import (
	"shop-web/common/dbutil"
	"gopkg.in/mgo.v2"
	"sync"
	"github.com/gin-gonic/gin"
	"shop-web/module/service/model"
	"shop-web/common/commonUtils"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type jobService struct {
	context *dbutil.MongoContext
	c       *mgo.Collection
	mutex   *sync.Mutex
}

func (js *jobService) SetContext() {
	js.context = dbutil.NewDBContext()
	js.c = js.context.DBCollection(Cname)
}

func (js *jobService) SaveService(c *gin.Context) {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	userId := c.GetString("userId")
	language := c.GetString("language")

	if userId == "" {
		Logger.Error("userId null")
		commonUtils.CreateErrorParams(c)
		return
	}

	var request model.ServiceJobResp
	if err := c.ShouldBindJSON(&request); err != nil {
		Logger.Error(err)
		commonUtils.CreateErrorParams(c)
		return
	}

	Logger.Info(request)
	location := model.GeoJson{
		Type:        "Point",
		Coordinates: request.Location,
	}

	serviceJob := model.ServiceJob{}
	serviceJob.Id = bson.NewObjectId()
	serviceJob.Category = request.Category
	serviceJob.UserId = userId
	serviceJob.Username = request.Username
	serviceJob.ContactPhone = request.ContactPhone
	serviceJob.State = request.State
	serviceJob.City = request.City
	serviceJob.Address = request.Address
	serviceJob.Area = request.Area
	serviceJob.Code = request.Code
	serviceJob.Require = request.Require
	serviceJob.Pics = request.Pics
	serviceJob.Title = request.Title
	serviceJob.Price = request.Price
	serviceJob.DescStr = request.DescStr
	serviceJob.Language = language
	serviceJob.Type = request.Type
	serviceJob.Status = 1
	serviceJob.UpdateAt = time.Now()
	serviceJob.CreateAt = time.Now()
	serviceJob.Location = location

	js.SetContext()
	defer js.context.Close()

	if err := js.c.Insert(&serviceJob); err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	commonUtils.CreateSuccess(c, serviceJob.Id)

}

func (js *jobService) FindService(c *gin.Context) {

}

var JobService = jobService{
	mutex: &sync.Mutex{},
}
