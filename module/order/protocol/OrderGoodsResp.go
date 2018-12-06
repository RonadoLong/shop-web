package protocol

type OrderGoodsResp struct {
	Id int64 `json:"id"`
	OrderId string `json:"orderId"`
	GoodsId int64 `json:"goodsId"`
	GoodsTitle string `json:"goodsTitle"`
	//单价
	GoodsPrice int `json:"goodsPrice"`
	//商品总价
	TotalPrice int `json:"totalPrice"`
	GoodsImage string `json:"goodsImage"`
	GoodsCount int `json:"goodsCount"`
	SkuId int64 `json:"skuId"`
	//属性名字
	SkuValues string `json:"skuValues"`
}