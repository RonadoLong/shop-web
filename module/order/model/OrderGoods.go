package model

type OrderGoods struct {

	Id int64
	OrderId string
	ProductId int64
	GoodsTitle string
	//单价
	GoodsPrice int
	//商品总价
	TotalPrice int
	GoodsImage string
	GoodsCount int
	//商品编号
	GoodsNumber string
	SkuId int64
	//属性名字
	SkuValues string

}
