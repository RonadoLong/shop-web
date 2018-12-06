package model

type UserAddress struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"  `
	ContactName string `json:"contactName"`
	Mobile      string `json:"mobile"`
	State       string `json:"state"`
	Address     string `json:"address" `
	PostalCode  string `json:"postalCode"`
}
