package service

import (
	"shop-web/common/protocol"
	"shop-web/common/redis"
	"strconv"
	"shop-web/common/dbutil"
	"shop-web/module/goods/model"
	"shop-web/common/commonUtils"
)



func (goodsService *goodsService) ReduceStock(reduceGoodsList []protocol.ReduceGoodsNumberReq) error {

	goodsService.Mutex.Lock()
	defer goodsService.Mutex.Unlock()

	var err error
	tx := dbutil.DB.Begin()

	//以订单ID 加个缓存锁 防止程序短时间重试 重复扣减库存 不用解锁 自己超时
	//todo

	//扣减每个商品的redis库存 扣减流水插入流水表
	for _, reduceGoodsCountsReq := range reduceGoodsList {

		skuId := reduceGoodsCountsReq.SkuId
		stockKey := redis.GOODS_STOCK_KEY + strconv.FormatInt(skuId, 10)
		stockLockKey := redis.GOODS_STOCK_LOCK_KEY + strconv.FormatInt(skuId, 10)
		reduceCount := -reduceGoodsCountsReq.GoodsCount
		reduceLockCount := reduceGoodsCountsReq.GoodsCount

		ret, err := redis.ReduceSkuStock(stockKey,stockLockKey,reduceCount,reduceLockCount)
		if err != nil {
			Logger.Error("ReduceSkuStock err", err)
			err = commonUtils.ErrStockNumberNotEnough
			break
		}

		switch ret.(type) {
		case int64:
			Logger.Error("Err StockNumber NotEnough")
			err = commonUtils.ErrStockNumberNotEnough
			break

		case []interface{}:
			//扣减成功 记录扣减流水
			r := ret.([]interface{})
			stockAftChange := int(r[0].(int64))
			stockLockAftChange := int(r[1].(int64))

			goodsStockFlow := model.GoodsStockFlow{}
			goodsStockFlow.SkuId = skuId
			goodsStockFlow.OrderId = reduceGoodsCountsReq.OrderId

			goodsStockFlow.LockStockAfter = stockLockAftChange
			goodsStockFlow.LockStockBefore = stockLockAftChange - reduceLockCount
			goodsStockFlow.LockStockChange = reduceLockCount

			goodsStockFlow.StockAfter = stockAftChange
			goodsStockFlow.StockBefore = stockAftChange - reduceCount
			goodsStockFlow.StockChange = reduceCount
			goodsStockFlow.CheckStatus = 0
			err = tx.Table("goods_stock_flow").Create(&goodsStockFlow).Error
			if err != nil {
				tx.Rollback()
				break
			}
		}

	}

	tx.Commit()
	if err != nil{
		return err
	}
	return nil
}


func (goodsService *goodsService) QueryStock(skuId int64) int {

	stockKey := redis.GOODS_STOCK_KEY + strconv.FormatInt(skuId, 10)
	stockLockKey := redis.GOODS_STOCK_LOCK_KEY + strconv.FormatInt(skuId, 10)

	var stockInRedis = 0
	ret, err := redis.RedisClient.Get(stockKey).Result()
	if err != nil{
		return 0
	}

	stockInRedis, err = strconv.Atoi(ret)
	//redis中为空 初始化一次
	if stockInRedis == 0{
		var goodsSku model.GoodsSku
		err := dbutil.DB.Table("goods_sku").Where("sku_id = ?",skuId).Find(&goodsSku).Error
		if err != nil{
			return 0
		}
		redis.SkuStockInit(stockKey,stockLockKey,goodsSku.Stock,goodsSku.LockStock)
		return goodsSku.Stock
	}

	return stockInRedis
}