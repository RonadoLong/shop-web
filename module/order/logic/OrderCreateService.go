package logic

import (
	"errors"
	"shop-web/common/commonUtils"
	"shop-web/common/protocol"
	model2 "shop-web/module/goods/model"
	"shop-web/module/goods/service"
	"shop-web/module/order/model"
	"strconv"

	"github.com/HaroldHoo/id_generator"
	"github.com/jinzhu/gorm"
)

func SaveNoConfirmOrder(confirmOrderReq *protocol.ConfirmOrderReq, tx *gorm.DB) (string, error) {
	var dataID = uint64(128)
	id, err := id_generator.NextId(dataID)
	if err != nil {
		return "", nil
	}

	totalIntegral := 0
	totalCount := 0
	totalPrice := 0
	orderIdStr := strconv.FormatUint(id, 10)

	for _, orderGoodsReq := range confirmOrderReq.OrderGoodsReqList {
		var product model2.Product

		err := tx.Table("product").
			Where("product_id = ? and status = 1 and stock >= ? ", orderGoodsReq.ProductId, orderGoodsReq.GoodsCount).Find(&product).Error
		if err != nil {
			logger.Error(err)
			return orderIdStr, errors.New("stock not enough")
		}

		if product.ProductId == 0 {
			logger.Error("stock not enough")
			return orderIdStr, errors.New("stock not enough")
		}

		err = tx.Exec("update product set stock = stock - ? where product_id = ?",
			orderGoodsReq.GoodsCount, product.ProductId).Error
		if err != nil {
			logger.Error(err)
			return orderIdStr, err
		}

		orderGoods := model.OrderGoods{}
		orderGoods.GoodsTitle = product.Title
		orderGoods.GoodsImage = product.GoodsImages
		orderGoods.GoodsPrice = product.MemberPrice
		orderGoods.GoodsCount = orderGoodsReq.GoodsCount
		orderGoods.ProductId = orderGoodsReq.ProductId
		//orderGoods.SkuId = orderGoodsReq.SkuId
		orderGoods.OrderId = orderIdStr
		//orderGoods.SkuValues = orderGoodsReq.SkuValues
		orderGoods.TotalPrice = product.MemberPrice * orderGoodsReq.GoodsCount
		totalCount += orderGoods.GoodsCount
		totalPrice += orderGoods.TotalPrice
		totalIntegral += product.Integral

		err = tx.Table("order_goods").Create(&orderGoods).Error
		if err != nil {
			return orderIdStr, err
		}

		err = tx.Exec("update shopping_cart set status = 0 where product_Id = ?", orderGoodsReq.ProductId).Error
		if err != nil {
			logger.Error(err)
		}
	}

	/** create order */
	orderInfo := model.OrderInfo{}
	orderInfo.OrderId = orderIdStr
	orderInfo.UserId = confirmOrderReq.UserId
	//orderInfo.Username = confirmOrderReq.Username
	orderInfo.GoodsCount = totalCount
	orderInfo.OrderAddress = confirmOrderReq.AddressId
	orderInfo.OrderIdentifier = strconv.FormatUint(id, 36)
	orderInfo.PayType = ""
	//orderInfo.Phone = confirmOrderReq.Phone
	/** 默认订单类型 */
	orderInfo.OrderType = 0
	/** 是否包邮 */
	orderInfo.IsPostFee = 1
	orderInfo.PostFee = 0

	if confirmOrderReq.CouponId != "" {
		orderInfo.CouponId = confirmOrderReq.CouponId
	}

	orderInfo.TotalAmount = totalPrice
	orderInfo.ReallyAmount = totalPrice
	orderInfo.TotalIntegral = totalIntegral
	orderInfo.OrderStatus = commonUtils.ORDERSTATUS_CONFIRMED
	orderInfo.PayStatus = commonUtils.PAYSTATUS_NO_PAY
	orderInfo.ShippingStatus = commonUtils.SHIPPINGSTATUS_NO_SHIP

	orderInfoErr := tx.Table("order_info").Create(&orderInfo).Error
	if orderInfoErr != nil {
		return orderIdStr, orderInfoErr
	}

	//err = tx.Exec("update user set integral = integral + ? where user_id = ?", totalIntegral, orderInfo.UserId).Error
	//if err != nil {
	//	return orderIdStr, err
	//}

	return orderIdStr, nil
}

func CallRemoteService(confirmOrderReq *protocol.ConfirmOrderReq, orderId string, tx *gorm.DB) error {

	var err error
	/** 扣库存 */
	reduceList := []protocol.ReduceGoodsNumberReq{}
	for _, orderGoodsReq := range confirmOrderReq.OrderGoodsReqList {
		reduce := protocol.ReduceGoodsNumberReq{}
		reduce.OrderId = orderId
		reduce.GoodsCount = orderGoodsReq.GoodsCount
		reduce.ProductId = orderGoodsReq.ProductId
		reduceList = append(reduceList, reduce)
	}

	err = service.GoodsService.ReduceStock(reduceList)
	if err != nil {
		return err
	}

	/** 更改订单状态 */
	err = tx.Table("order_info").
		Where("order_id = ?", orderId).Updates(model.OrderInfo{OrderStatus: commonUtils.ORDERSTATUS_CONFIRMED}).Error
	if err != nil {
		return err
	}

	return nil
}
