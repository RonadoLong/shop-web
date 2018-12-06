package service

import (
	"github.com/gin-gonic/gin"
	"shop-web/common/dbutil"
	"shop-web/module/goods/model"
	"shop-web/common/commonUtils"
	"shop-web/module/goods/protocol"
)



func (*goodsService)FindGoodsNavList(c *gin.Context)   {

	var goodsNavList []model.GoodsNav
	if err := dbutil.DB.Table("goods_nav").Where("status = 1").Order("sort asc").Find(&goodsNavList).Error; err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	var goodsNavRespList = []protocol.GoodsNavResp{}
	for _,goodsNav := range goodsNavList {
		goodsNavResp := protocol.GoodsNavResp{}
		goodsNavResp.ClassId = goodsNav.ClassId
		if c.GetString("language") == "EN" {
			goodsNavResp.Title = goodsNav.EnTitle
		}else {
			goodsNavResp.Title = goodsNav.Title
		}
		goodsNavRespList = append(goodsNavRespList,goodsNavResp)
	}

	commonUtils.CreateSuccess(c,goodsNavRespList)
}
