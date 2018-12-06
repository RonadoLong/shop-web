package cron

import (
	"github.com/robfig/cron"
	"shop-web/common/dbutil"
	"shop-web/common/redis"
	stock "shop-web/module/goods/model"
	"os"
	"shop-web/common/log"
)

var Logger = log.NewLogger(os.Stdout)

func InitCron(stockCron string)  {
	Logger.Info("init cron")
	c := cron.New()
	//c.AddFunc(stockCron, StockSimpleJob)
	c.Start()
}

type Result struct {
	StockAfter int64
	LockStockAfter int64
}


/**
	清算库存， 更新数据库库存
 */
func StockSimpleJob() {

	Logger.Info("begin stock job")

	var err error
	var goodsStockFlowList []stock.GoodsStockFlow
	err = dbutil.DB.Table("goods_stock_flow").Where("`check_status` = 0").Find(&goodsStockFlowList).Error
	Logger.Info(goodsStockFlowList)
	if err == nil && len(goodsStockFlowList) > 0{
		for _,goodsStockFlow := range goodsStockFlowList {
			skuId := goodsStockFlow.SkuId
			id := goodsStockFlow.Id
			Logger.Infof("skuId %d === id %d", skuId, goodsStockFlow.Id)

			//redis 过期 是否存在 orderId
			ret :=  redis.RedisClient.Exists("order:" + goodsStockFlow.OrderId).Val()
			Logger.Info(ret)
			if int(ret) > 0 {
				return
			}

			result := Result{}
			rows,err := dbutil.DB.Table("goods_stock_flow").Select("min(`stock_after`) as stockAfter, max(`lock_stock_after`) as lockStockAfter").Where("`sku_id` = ? and `check_status` = 0", skuId).Rows()
			for rows.Next(){
				rows.Scan(&result.StockAfter, &result.LockStockAfter)
			}
			Logger.Info(result)

			if err != nil{
				Logger.Error(err)
				return
			}

			err = dbutil.DB.Table("goods_sku").Where("`sku_id` = ?", skuId).Updates(stock.GoodsSku{Stock: int(result.StockAfter), LockStock: int(result.LockStockAfter)}).Error
			if err != nil{
				Logger.Error(err)
				return
			}

			err = dbutil.DB.Table("goods_stock_flow").Where("`id` = ?", id).Updates(stock.GoodsStockFlow{CheckStatus:1}).Error
			if err != nil{
				Logger.Error(err)
				return
			}

		}

	}

	Logger.Info("finish stock job")
}