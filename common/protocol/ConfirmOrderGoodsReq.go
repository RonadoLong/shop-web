package protocol

type ConfirmOrderGoodsReq struct {

	ProductId int64
	GoodsTitle string
	GoodsPrice int
	TotalPrice int
	GoodsImage string
	GoodsCount int
	//SkuId int64
	//SkuValues string
}
