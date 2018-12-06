package protocol

type OrderInfoResp struct {
	OrderId      string `json:"orderId"`
	UserId       string `json:"userId"`
	Username     string `json:"username"`
	MerchantId   string
	OrderAddress string `json:"orderAddress"`
	Phone        string `json:"phone"`
	//是否包邮
	IsPostFee int `json:"isPostFee"`
	//邮费
	PostFee int `json:"postFee"`
	//商品数量
	GoodsCount int `json:"goodsCount"`
	//总金额
	TotalAmount int `json:"totalAmount"`
	//实际支付金额
	ReallyAmount int `json:"reallyAmount"`
	//订单编码
	OrderIdentifier string `json:"orderIdentifier"`
	//买家是否已经评价
	BuyerRate int `json:"buyerRate"`
	//支付类型 1 pay pal  2 信用卡
	PayType string `json:"payType"`

	//状态：1未确认 2已确认 3退款 4交易成功(已收货) 5交易关闭 6无效
	OrderStatus int `json:"orderStatus"`

	//发货状态 0未发货 1已发货 2已收货
	ShippingStatus int `json:"shippingStatus"`

	//支付状态 0未支付 1已支付
	PayStatus int `json:"payStatus"`

	Status int `json:"status"`

	OrderGoodsRespList []OrderGoodsResp `json:"orderGoodsRespList"`
}
