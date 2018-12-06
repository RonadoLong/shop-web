package logic

import (
	"sync"
	"shop-web/common/dbutil"
	"gopkg.in/mgo.v2"
	"shop-web/module/service/model"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type servicePayFlowService struct {
	mutex *sync.Mutex
	context *dbutil.MongoContext
	c       *mgo.Collection
}

var (
	ServicePayFlowService = servicePayFlowService{
		mutex: &sync.Mutex{},
	}

	ServicePayFlowCname = "service_order_flow"
)

func (ad *servicePayFlowService) SetContext() {
	ad.context = dbutil.NewDBContext()
	ad.c = ad.context.DBCollection(ServicePayFlowCname)
}

func (ad *servicePayFlowService)UpdateStatus(parmas *model.ServicePayFlowProtocol){

	ex_time := time.Now().AddDate(0,0, parmas.ExTime)
	servicePayFlow := model.ServicePayFlow{
		Id: bson.NewObjectId(),
		ServiceId: parmas.ServiceId,
		ExTime: ex_time,
		Price: parmas.Price,
		Status: 1,
		UpdateAt: time.Now(),
		CreateAt: time.Now(),
	}

	ad.SetContext()
	defer ad.context.Close()

	if err := ad.c.Insert(&servicePayFlow); err != nil{
		Logger.Error(err)
		return
	}

	AdtService.UpdateServiceStatus(parmas.ServiceId)

	Logger.Info("update ServiceId f")
}