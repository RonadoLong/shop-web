package controller

import (
	"shop-web/module/service/logic"

	"github.com/gin-gonic/gin"
)

var ServiceRest = serviceRest{}

type serviceRest struct {
}

func (*serviceRest) FindServiceCategoryList(c *gin.Context) {
	logic.ServiceClassService.FindClassList(c)
}

func (*serviceRest) FindAreaListByParentId(c *gin.Context) {
	logic.FindAreatListByParentId(c)
}

func (*serviceRest) AddServiceRoom(c *gin.Context) {
	logic.AdtService.SaveService(c)
}

func (*serviceRest) FindServiceRoomByCategory(c *gin.Context) {
	logic.AdtService.FindServiceList(c)
}

func (*serviceRest) FindServiceByCategory(c *gin.Context) {
	logic.AdtService.FindServiceList(c)
}

func (*serviceRest) SaveServiceJob(c *gin.Context) {
	logic.JobService.SaveService(c)
}

func (*serviceRest) FindSelfService(c *gin.Context)  {
	logic.AdtService.FindSelfService(c)
}

func (*serviceRest) FindServicePaymentList(c *gin.Context) {
	logic.AdtService.FindServicePaymentList(c)
}
