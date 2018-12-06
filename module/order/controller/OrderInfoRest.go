package controller

import (
	"shop-web/module/order/logic"

	"github.com/gin-gonic/gin"
)

var OrderInfoRest = orderInfoRest{}

type orderInfoRest struct {
}

func (*orderInfoRest) CreateOrder(c *gin.Context) {
	logic.OrderInfoService.CreateOrder(c)
}

func (*orderInfoRest) FindOrderInfoList(c *gin.Context) {
	logic.OrderInfoService.FindOrderInfoList(c)
}

func (*orderInfoRest) CancelOrder(c *gin.Context) {
	logic.OrderInfoService.CancelOrder(c)
}

func (*orderInfoRest) ReceiveIPN(c *gin.Context) {
	logic.OrderInfoService.ReceiveIPN(c)
}
