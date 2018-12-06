package logic

import (
	"shop-web/module/order/model"
	"shop-web/common/commonUtils"
)

//1 未支付 2 待发货/已支付 3已发货 4交易成功 5交易关闭 6无效
func ReturnStatus(info model.OrderInfo) int {

	if info.PayStatus == commonUtils.PAYSTATUS_NO_PAY {
		return 1
	}

	if info.PayStatus == commonUtils.PAYSTATUS_PAIED && info.ShippingStatus == commonUtils.SHIPPINGSTATUS_NO_SHIP {
		return 2
	}

	if info.PayStatus == commonUtils.PAYSTATUS_PAIED && info.ShippingStatus == commonUtils.SHIPPINGSTATUS_SHIPED {
		return 3
	}

	if info.PayStatus == commonUtils.PAYSTATUS_PAIED && info.OrderStatus == commonUtils.ORDERSTATUS_SUCCESSFULL{
		return 4
	}

	if info.PayStatus == commonUtils.PAYSTATUS_PAIED && info.OrderStatus == commonUtils.ORDERSTATUS_ClOSE{
		return 5
	}

	if info.PayStatus == commonUtils.PAYSTATUS_PAIED && info.OrderStatus == commonUtils.ORDERSTATUS_INVAILD{
		return 6
	}

	return 6
}