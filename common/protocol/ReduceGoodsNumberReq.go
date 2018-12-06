package protocol

type ReduceGoodsNumberReq struct {
	OrderId string
	ProductId int64
	SkuId int64
	GoodsCount int
} 