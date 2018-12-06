package dao

import (
	"shop-web/common/dbutil"
	"shop-web/module/order/model"
	"shop-web/module/order/protocol"
	"shop-web/common/commonUtils"
)

//	PayType int

//	//状态：1未确认 2已确认 3退款 4交易成功(已收货) 5交易关闭 6无效
//	OrderStatus int

//	//发货状态 1未发货 2已发货 3已收货
//	ShippingStatus int

//	//支付状态 1未支付 2已支付
//	PayStatus int

func FindOrderInfoList(offset int, pageSize int, dataType int, userId string) ([]model.OrderInfo, error){
	var orderInfoList []model.OrderInfo
	var err error
	//全部
	if dataType == 0{
		err = dbutil.DB.Table("order_info").
			Where("`user_id` = ? ", userId).Offset(offset).Limit(pageSize).Find(&orderInfoList).Error
		return orderInfoList, err
	}

	//等待支付
	if dataType == 1{
		err = dbutil.DB.Table("order_info").
			Where("user_id = ? and pay_status = 1 and order_status = 2", userId).Offset(offset).Limit(pageSize).
			Find(&orderInfoList).Error
		return orderInfoList, err
	}

	//待发货
	if dataType == 2{
		err = dbutil.DB.Table("order_info").
			Where("user_id = ? and shipping_status = 1 and pay_status = 2", userId).Offset(offset).Limit(pageSize).
			Find(&orderInfoList).Error
		return orderInfoList, err
	}

	//已发货
	if dataType == 3{
		err = dbutil.DB.Table("order_info").
			Where("user_id = ? and shipping_status = 2 and pay_status = 2", userId).Offset(offset).Limit(pageSize).
			Find(&orderInfoList).Error
		return orderInfoList, err
	}

	return nil, err
}


func FindOrderInfoCounts(dataType int, userId string) (int, error) {

	var err error
	count := 0

	//全部
	if dataType == 0{
		err = dbutil.DB.Table("order_info").
			Where("`user_id` = ? ", userId).Count(&count).Error
		return count, err
	}

	//等待支付
	if dataType == 1{
		err = dbutil.DB.Table("order_info").
			Where("user_id = ? and pay_status = 1 and order_status = 2", userId).Count(&count).Error
		return count, err
	}

	//待发货
	if dataType == 2{
		err = dbutil.DB.Table("order_info").
			Where("`user_id` = ? and `shipping_status` = 1 and `pay_status` = 2", userId).Count(&count).Error
		return count, err
	}

	//已发货
	if dataType == 3{
		err = dbutil.DB.Table("order_info").
			Where("user_id = ? and shipping_status = 2 and pay_status = 2", userId).Count(&count).Error
		return count, err
	}

	//待评价
	if dataType == 4{
		err = dbutil.DB.Table("order_info").
			Where("`user_id` = ? and `shipping_status` = 3 and `order_status` = 4", userId).Count(&count).Error
		return count, err
	}

	return count, nil
}

func FindOrderGoodsByOrderId(orderId string) ([]protocol.OrderGoodsResp, error) {

	var orderGoodsList []protocol.OrderGoodsResp
	if err := dbutil.DB.Table("order_goods").
		Where("order_id = ?", orderId).Find(&orderGoodsList).Error; err != nil {
		return nil, err
	}
	return orderGoodsList, nil
}

func CancelOrder(orderId string, userId string) error {
	if err := dbutil.DB.Exec(" update `order_info` set order_status = ? where `user_id` = ? and `order_id` = ?", commonUtils.ORDERSTATUS_ClOSE, userId, orderId).Error; err != nil {
		return err
	}
	return nil
}