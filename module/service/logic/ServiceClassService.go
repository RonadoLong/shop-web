package logic

import (
	"shop-web/common/dbutil"
	"gopkg.in/mgo.v2"
	"sync"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"shop-web/common/commonUtils"
	"shop-web/module/service/model"
)

const COLLECTION_NAME  =  "service_class"

type serviceClassService struct {
	context *dbutil.MongoContext
	c       *mgo.Collection
	mutex   *sync.Mutex
}

func (ad *serviceClassService) SetContext() {
	ad.context = dbutil.NewDBContext()
	ad.c = ad.context.DBCollection(COLLECTION_NAME)
}

var ServiceClassService = serviceClassService {
	mutex: &sync.Mutex{},
}

func (serviceClassService *serviceClassService) FindClassList(c * gin.Context){
	serviceClassService.SetContext()
	defer  serviceClassService.context.Close()
	query := bson.M{
		"status": 1,
	}
	total, err := serviceClassService.c.Find(query).Count()
	language := c.GetString("language")

	if total == 0 {
		Logger.Info("serviceClassService no find")
		commonUtils.CreateNotContent(c)
		return
	}

	var serviceClassList []model.ServiceCLass
	err = ServiceClassService.c.Find(query).All(&serviceClassList)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateErrorRequest(c)
		return
	}

	if language == "EN" {
		serviceClassListResp := make([]model.ServiceCLass, 0)
		for _, serviceClass := range serviceClassList {
			serviceClass.Name = serviceClass.EnName
			serviceClassListResp = append(serviceClassListResp, serviceClass)
		}
		commonUtils.CreateSuccess(c, serviceClassListResp)
		return
	}

	commonUtils.CreateSuccess(c, serviceClassList)
}

