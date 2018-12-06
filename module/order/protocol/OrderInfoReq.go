package protocol

import "time"

type OrderInfoReq struct {
	PageNum  int
	PageSize int

	//0 全部 1 等待支付 2 待发货
	DataType int
}

type ReceiviReq struct {
	Amount      string
	Currency    string
	Id          string
	Rmb_amount  string
	Reference   string
	Note        string
	Status      string
	Sys_reserve Sys_reserve
	Time        time.Time
	Verify_sign string
}

type Sys_reserve struct {
	Vendor_id string
}
