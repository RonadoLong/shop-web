package redis

import (
	"os"
	"shop-web/common/commonUtils"
	"shop-web/common/dbutil"
	"shop-web/common/log"
	"shop-web/module/goods/model"
	model2 "shop-web/module/order/model"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	StockAfter     int64
	LockStockAfter int64
}

var Logger = log.NewLogger(os.Stdout)

/**
订单过期事件 处理redis中的库存
*/
func PSubscribeEventHandel() {

	sub := RedisClient.PSubscribe("__keyevent@0__:expired")
	sub.ReceiveTimeout(time.Second * 110)

	for message := range sub.Channel() {
		Logger.Info("__keyevent === expired , ", message.Payload)
		orderId := strings.Split(message.Payload, ":")[1]

		if orderId != "" && !strings.Contains(orderId, ".") {
			go resetGoodsCount(orderId)
		}
		Logger.Info("finish, ", message.Payload)
	}

}

func resetGoodsCount(orderId string) {
	var (
		orderGoodsList []model2.OrderGoods
		err            error
	)
	tx := dbutil.DB.Begin()
	defer tx.Commit()

	err = tx.Table("order_goods").Where("order_id = ?", orderId).Find(&orderGoodsList).Error
	if err != nil || len(orderGoodsList) == 0 {
		goto ERR
	}

	/** 更新订单为失效 */
	err = tx.Table("order_info").Where("order_id = ? ", orderId).Updates(model2.OrderInfo{OrderStatus: commonUtils.ORDERSTATUS_INVAILD}).Error
	if err != nil {
		goto ERR
	}

	for _, orderGoods := range orderGoodsList {
		productId := orderGoods.ProductId
		stock := orderGoods.GoodsCount

		/** 恢复库存 */
		err := tx.Exec("update product set stock = stock + ? where product_id = ? ", stock, productId).Error
		if err != nil {
			goto ERR
		}
	}
ERR:
	Logger.Error(err)
	tx.Rollback()
}

func checkOrder(orderId string) {
	var err error
	var orderGoodsList []model2.OrderGoods
	err = dbutil.DB.Table("order_goods").Where("order_id = ?", orderId).Find(&orderGoodsList).Error
	if err != nil || len(orderGoodsList) == 0 {
		Logger.Error(err)
		return
	}

	for _, orderGoods := range orderGoodsList {
		go func() {
			skuId := orderGoods.SkuId
			stockChange := orderGoods.GoodsCount
			Logger.Info("skuId ==== ", skuId)

			/** 更新流水为检查过 避免重复更新*/
			err = dbutil.DB.Exec("update goods_stock_flow set check_status = 1 where sku_id = ? ", skuId).Error
			if err != nil {
				Logger.Error(err)
				return
			}

			/** 更新订单为失效 */
			err = dbutil.DB.Table("order_info").Where("order_id = ? ", orderId).Updates(model2.OrderInfo{OrderStatus: commonUtils.ORDERSTATUS_INVAILD}).Error
			if err != nil {
				Logger.Error(err)
			}

			/** 恢复库存 */
			stockKey := GOODS_STOCK_KEY + strconv.FormatInt(skuId, 10)
			stockLockKey := GOODS_STOCK_LOCK_KEY + strconv.FormatInt(skuId, 10)
			Logger.Info(stockKey, stockLockKey)

			ret, err := ReduceSkuStock(stockKey, stockLockKey, stockChange, -stockChange)
			Logger.Info(ret)
			if err != nil {
				Logger.Error(err)
			}
		}()
	}
}

func check1(orderId string) {
	var err error
	var goodsStockFlowList []model.GoodsStockFlow
	err = dbutil.DB.Table("goods_stock_flow").Where("check_status = 0 and order_id = ?", orderId).Find(&goodsStockFlowList).Error

	Logger.Info(goodsStockFlowList)
	if err == nil && len(goodsStockFlowList) > 0 {
		tx := dbutil.DB.Begin()
		for _, goodsStockFlow := range goodsStockFlowList {

			skuId := goodsStockFlow.SkuId
			id := goodsStockFlow.Id
			result := Result{}
			rows, err := tx.Table("goods_stock_flow").Select("min(stock_after) as stockAfter, max(lock_stock_after) as lockStockAfter").Where("sku_id = ?", skuId).Rows()
			for rows.Next() {
				rows.Scan(&result.StockAfter, &result.LockStockAfter)
			}

			err = tx.Table("goods_sku").Where("sku_id = ? ", skuId).Updates(model.GoodsSku{Stock: int(result.StockAfter), LockStock: int(result.LockStockAfter)}).Error
			if err != nil {
				Logger.Error(err)
				tx.Rollback()
				break
			}

			err = tx.Table("goods_stock_flow").Where("id = ?", id).Updates(model.GoodsStockFlow{CheckStatus: 1}).Error
			if err != nil {
				Logger.Error(err)
				tx.Rollback()
				break
			}

			err = tx.Table("order_info").Where("order_id = ? ", orderId).Updates(model2.OrderInfo{OrderStatus: commonUtils.ORDERSTATUS_INVAILD}).Error
			if err != nil {
				Logger.Error(err)
				tx.Rollback()
				break
			}

			stockKey := GOODS_STOCK_KEY + strconv.FormatInt(skuId, 10)
			stockLockKey := GOODS_STOCK_LOCK_KEY + strconv.FormatInt(skuId, 10)
			_, err = SkuStockInit(stockKey, stockLockKey, result.StockAfter, result.LockStockAfter)
			if err != nil {
				break
			}
		}

		tx.Commit()
	}
}
