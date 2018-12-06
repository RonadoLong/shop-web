package model

import "time"

type OrderInfo struct {
	OrderId string
	UserId string
	Username string
	MerchantId string `gorm:"default:'galeone'"`
	OrderAddress string
	OrderType int
	//是否包邮
	IsPostFee int
	//邮费。精确到2位小数;单位:元。如:200.07，表示:200元7分
	PostFee int
	//优惠券
	CouponId string `gorm:"default:'galeone'"`
	//优惠券金额'
	CouponPaid int `gorm:"default:'galeone'"`
	//商品数量
	GoodsCount int

	TotalIntegral int
	//总金额
	TotalAmount int
	//实际支付金额
	ReallyAmount int
	//订单编码
	OrderIdentifier string `gorm:"default:'galeone'"`
	//物流名称
	ShippingName string `gorm:"default:'galeone'"`
	//物流单号
	ShippingCode string `gorm:"default:'galeone'"`
	//买家留言
	BuyerMsg string `gorm:"default:'galeone'"`
	//付款时间
	PaymentTime time.Time `gorm:"default:'galeone'"`
	//订单结算总价
	TotalSettlementPrice int `gorm:"default:'galeone'"`
	//买家是否已经评价
	BuyerRate int `gorm:"default:'galeone'"`
	//支付类型 1 pay pal  2 信用卡
	PayType string
	//状态：1未确认 2已确认 3退款 4交易成功(已收货) 5交易关闭 6无效
	OrderStatus int
	//发货状态 1未发货 2已发货 3已收货
	ShippingStatus int
	//支付状态 1未支付 2支付中 3已支付
	PayStatus int
	////订单创建时间
	//CreateTime time.Time
	////订单更新时间
	//UpdateTime time.Time
}

