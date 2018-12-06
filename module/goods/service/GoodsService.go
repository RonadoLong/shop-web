package service

import (
	"github.com/gin-gonic/gin"
	"shop-web/common/log"
	"os"
	"shop-web/common/commonUtils"
	"shop-web/module/goods/dao"
	"shop-web/common/protocol"
	"sync"
	"shop-web/common/dbutil"
	"strconv"
	protocol2 "shop-web/module/goods/protocol"
)

var Logger = log.NewLogger(os.Stdout)

var GoodsService = goodsService{
	Mutex: &sync.Mutex{},
}

type goodsService struct {
	Mutex *sync.Mutex
}

func (*goodsService)FindGoodsList(c *gin.Context)  {
	offset,pageSize:= commonUtils.GetOffset(c)
	category,_ := strconv.Atoi(c.Param("category"))

	total,err := dao.FindGoodsCount(category)
	if err != nil{
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	if offset >= total{
		commonUtils.CreateNotContent(c)
		return
	}

	goodsList,err := dao.FindGoodsListByOffset(offset, pageSize, category)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}

	goodsRespList := []protocol2.GoodsListResp{}
	for _,goods := range goodsList {
		goodsResp,err := commonUtils.TransferValues(goods, c.GetString("language"))
		if err != nil {
			Logger.Error(err)
		}
		Logger.Info(goodsResp)
		goodsRespList = append(goodsRespList, goodsResp)
	}

	ret := gin.H{
		"total": total,
		"content": goodsRespList,
	}

	commonUtils.CreateSuccess(c, ret)
}

//func (goodsService *goodsService) FindGoodsByGoodsList(orderGoodsList []protocol.ConfirmOrderGoodsReq, language string) ([]protocol.ConfirmOrderGoodsReq, error) {
//
//	confirmGoodsList := []protocol.ConfirmOrderGoodsReq{}
//
//	for _,orderGoods := range orderGoodsList {
//		goodsSku,err := dao.FindGoodsSkuByGoodsIdAndSkuId(orderGoods.SkuId)
//
//		if err != nil && goodsSku.SkuId == 0 {
//			err = commonUtils.ErrSkuIdAndGoodsIdWrong
//			return nil, err
//		}
//
//		if orderGoods.GoodsCount > goodsSku.Stock {
//			err := commonUtils.ErrStockNumberNotEnough
//			Logger.Error(err)
//			return nil, err
//		}
//
//		orderGoods.TotalPrice = goodsSku.MemberPrice * orderGoods.GoodsCount
//		goods,err := dao.FindGoodsByGoodsId(goodsSku.GoodsId)
//		if err != nil || goods.ProductId == 0{
//			err = commonUtils.ErrStockNumberNotEnough
//			Logger.Error(err)
//			return nil, err
//		}
//
//		orderGoods.GoodsPrice = goodsSku.MemberPrice
//		orderGoods.GoodsImage = goods.GoodsImages
//		if language == "EN" {
//			orderGoods.GoodsTitle = goods.EnTitle
//			//todo
//			orderGoods.SkuValues = "测试"
//		}else {
//			orderGoods.GoodsTitle = goods.Title
//		}
//
//		confirmGoodsList = append(confirmGoodsList, orderGoods)
//	}
//
//	return confirmGoodsList, nil
//}
//

func (goodsService *goodsService) ReduceGoodsCounts(reduceGoodsList []protocol.ReduceGoodsNumberReq) error {

	goodsService.Mutex.Lock()
	defer goodsService.Mutex.Unlock()

	tx := dbutil.DB.Begin()
	for _, reduceGoodsCountsReq := range reduceGoodsList {
		skuId := reduceGoodsCountsReq.SkuId
		goodsCount := reduceGoodsCountsReq.GoodsCount
		if err := tx.Exec(" update `goods_sku` set `stock` = `stock` - ? where (`stock` > ? or `stock` = ?) and sku_id = ? ", goodsCount,goodsCount,goodsCount,skuId).Error; err != nil{
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (goodsService *goodsService)FindGoodsDetail(c *gin.Context){

	goodsId,_ := strconv.ParseInt(c.Param("goodsId"),10,64)
	if goodsId == 0{
		 commonUtils.CreateErrorParams(c)
		 return
	}

	goods,err := dao.FindGoodsByGoodsId(goodsId)
	if err != nil{
		Logger.Error(err)
	}

	goodsResp,err := commonUtils.TransferValues(goods,c.GetString("language"))
	commonUtils.CreateSuccess(c,goodsResp)
}

func (goodsService *goodsService) FindGoodsAllList(c *gin.Context) {
	offset,pageSize:= commonUtils.GetOffset(c)
	total,err := dao.FindGoodsAllCount()
	if err != nil{
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}
	if offset >= total{
		commonUtils.CreateNotContent(c)
		return
	}
	goodsList,err := dao.FindAllGoodsListByOffset(offset, pageSize)
	if err != nil {
		Logger.Error(err)
		commonUtils.CreateError(c)
		return
	}
	Logger.Info(goodsList)

	goodsRespList := make([]protocol2.GoodsListResp, 0)
	Logger.Info(goodsRespList)

	for _,goods := range goodsList {
		goodsResp,err := commonUtils.TransferValues(goods, c.GetString("language"))
		if err != nil {
			Logger.Error(err)
		}else {
			goodsRespList = append(goodsRespList, goodsResp)
			Logger.Info(len(goodsRespList))
		}
	}

	ret := gin.H{
		"total": total,
		"content": goodsRespList,
	}
	commonUtils.CreateSuccess(c, ret)
}

