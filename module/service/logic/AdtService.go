package logic

import (
	"shop-web/common/commonUtils"
	"shop-web/common/dbutil"
	"shop-web/module/service/model"
	"sync"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strconv"
)

type adtService struct {
	context *dbutil.MongoContext
	c       *mgo.Collection
	mutex   *sync.Mutex
}

var AdtService = adtService{
	mutex: &sync.Mutex{},
}

func (ad *adtService) SetContext() {
	ad.context = dbutil.NewDBContext()
	ad.c = ad.context.DBCollection(Cname)
}

func (ad *adtService) SaveService(c *gin.Context) {
	ad.mutex.Lock()
	defer ad.mutex.Unlock()

	userId := c.GetString("userId")
	language := c.GetString("language")

	var request map[string]interface{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateErrorParams(c)
		return
	}

	Logger.Info(userId, language)
	Logger.Info(request)
	request["userId"] = userId
	request["language"] = language
	status := request["status"]
	if status == nil{
		request["status"] = 0
	}

	request["update_at"] = time.Now()
	request["create_at"] = time.Now()

	request["_id"] = bson.NewObjectId()
	ok := request["location"]
	Logger.Info(ok)
	ad.SetContext()
	defer ad.context.Close()
	err = ad.c.Insert(&request)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	commonUtils.CreateSuccess(c, request["_id"])
}

func (ad *adtService) FindServiceList(c *gin.Context) {
	serviceReq := model.ServiceReq{}
	if err := c.ShouldBindJSON(&serviceReq); err != nil {
		Logger.Info(err)
		commonUtils.CreateErrorParams(c)
		return
	}

	Logger.Info("params === ", serviceReq)
	ad.FindServiceByType(c, &serviceReq)

}

//func (ad *adtService) FindServiceByNearby(c *gin.Context, serviceReq *model.ServiceReq) {
//	//如果要计算公里，可设置 distanceMultiplier: 6371
//	//如果要计算英里，则需要把6371换成3959
//
//	ad.SetContext()
//	defer ad.context.Close()
//
//	language := c.GetString("language")
//
//	query := bson.M{
//		"language": language,
//		"status":   1,
//	}
//
//	if serviceReq.Category != "" {
//		query["category"] = serviceReq.Category
//	}
//
//	if len(serviceReq.RoomType) != 0 {
//		query["roomtype"] = serviceReq.RoomType
//	}
//
//	Logger.Info("query === ", query)
//	pageNum, pageSize := commonUtils.GetOffset(c)
//	pipe := ad.c.Pipe([]bson.M{
//		{
//			"$geoNear": bson.M{
//				"near": bson.M{
//					"type":        "Point",
//					//"coordinates": []float64{serviceReq.Longitude, serviceReq.Latitude},
//				},
//				"query":              query,
//				"distanceMultiplier": 0.001, //km
//				"spherical":          true,  //代表弧度
//				"distanceField":      "distance",
//				"maxDistance":        4 * 1000,
//			},
//		},
//		{
//			"$sort": bson.M{
//				"distance": 1,
//				"createat": -1,
//			},
//		},
//		{
//			"$skip": pageNum,
//		},
//		{
//			"$limit": pageSize,
//		},
//	})
//
//	var resList = make([]map[string]interface{}, pageSize)
//
//	var result = make(map[string]interface{})
//
//	iter := pipe.Iter()
//	for iter.Next(&result) {
//		resList = append(resList, result)
//	}
//
//	if err := iter.Close(); err != nil {
//		Logger.Error(err)
//		return
//	}
//
//	if len(resList) == 0 {
//		commonUtils.CreateNotContent(c)
//		return
//	}
//
//	commonUtils.CreateSuccessByList(c, len(resList), resList)
//
//}

func (ad *adtService) FindServiceByType(c *gin.Context, serviceReq *model.ServiceReq) {
	//language := c.GetString("language")

	ad.SetContext()
	defer ad.context.Close()

	var count = 0

	query := bson.M{
		"status":   1,
	}

	if len(serviceReq.State) != 0 {
		query["state"] = serviceReq.State
	}

	if serviceReq.Category != "" {
		query["classId"] = serviceReq.Category
	}

	if len(serviceReq.RoomType) != 0 {
		query["roomType"] = serviceReq.RoomType
	}

	Logger.Info("query === ", query)

	count, err := ad.c.Find(query).Count()
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	pageNum, pageSize := commonUtils.GetOffset(c)
	Logger.Info("query === ", count, pageNum, pageSize)
	if count <= pageNum {
		commonUtils.CreateNotContent(c)
		return
	}
	var resList = make([]map[string]interface{}, 10)
	Logger.Info("query === ", query)
	err = ad.c.Find(query).Skip(pageNum).Limit(pageSize).Sort("-createat").All(&resList)
	if err != nil {
		Logger.Error(err)
		return
	}
	commonUtils.CreateSuccessByList (c, len(resList), resList)
}

func (ad *adtService) FindSelfService(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == ""{
		 commonUtils.CreateErrorParams(c)
		return
	}
	ad.SetContext()
	defer ad.context.Close()

	query := bson.M{
		"userId": userId,
	}
	status, _ := strconv.Atoi(c.Param("status"))
	if status != 6{
		query["status"] = status
	}

	count, err := ad.c.Find(query).Count()
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	if count == 0 {
		commonUtils.CreateNotContent(c)
		return
	}

	pageNum, pageSize := commonUtils.GetOffset(c)
	Logger.Info("query === ", count, pageNum, pageSize)
	if count <= pageNum {
		commonUtils.CreateNotContent(c)
		return
	}

	var resList = make([]map[string]interface{}, pageSize)
	Logger.Info("query === ", query)
	err = ad.c.Find(query).Skip(pageNum).Limit(pageSize).Sort("-create_at").All(&resList)
	if err != nil {
		Logger.Error(err)
		return
	}
	commonUtils.CreateSuccessByList (c, len(resList), resList)

}

/**
	get logic category
*/
func FindServiceCategoryList(c *gin.Context) {

	var categoryList []model.ServiceCategory
	err := dbutil.DB.Table("service_category").Where("status = 1").Find(&categoryList).Error
	if err != nil {
		Logger.Error(err.Error())
		commonUtils.CreateError(c)
		return
	}
	commonUtils.CreateSuccess(c, categoryList)
}


func (*adtService)FindServicePaymentList(c *gin.Context) {
	var paymentList []model.ServicePaymentSetting
	err := dbutil.DB.Table("service_payment_setting").
		Where("status = 1").Find(&paymentList).Error
	if err != nil {
		Logger.Error(err.Error())
		commonUtils.CreateError(c)
		return
	}
	commonUtils.CreateSuccess(c, paymentList)
}


func (ad *adtService)UpdateServiceStatus(serviceId string) {
	if serviceId == "" {
		return
	}

	ad.SetContext()
	defer ad.context.Close()

	selector := bson.M{"_id": bson.ObjectIdHex(serviceId)}
	data := bson.M{"$set": bson.M{"status": 1}}

	if err := ad.c.Update(selector, data); err != nil {
		Logger.Error(err)
	}

}

