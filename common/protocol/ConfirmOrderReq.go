package protocol

type ConfirmOrderReq struct {

	UserId string

	AddressId string

	Username string

	Phone string

	//默认是paypal
	PayType int

	// 优惠券ID
	CouponId string

	// 商品列表
	OrderGoodsReqList []ConfirmOrderGoodsReq
}
