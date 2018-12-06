package commonModel

import "time"

type Model struct {
	CreateTime time.Time `json:"createTime" `
	UpdateTime time.Time `json:"updateTime" `
}
