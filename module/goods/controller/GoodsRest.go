package controller

import (
	"github.com/gin-gonic/gin"
	"shop-web/module/goods/service"
)

var GoodsRest = goodsRest{}

type goodsRest struct {

}

func (*goodsRest)FindGoodsList(c *gin.Context)  {
	service.GoodsService.FindGoodsList(c)
}
func (*goodsRest)FindGoodsAllList(c *gin.Context)   {
	service.GoodsService.FindGoodsAllList(c)
}

func (*goodsRest)FindGoodsNavList(c *gin.Context) {
	service.GoodsService.FindGoodsNavList(c)
}

func (*goodsRest)FindGoodsDetail(c *gin.Context){
	service.GoodsService.FindGoodsDetail(c)
}