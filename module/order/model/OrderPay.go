package model

type OrderPay struct {
	//支付编号
	PayId string
	OrderId string
	PayAmount float64
	IsPaid string
}

