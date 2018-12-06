package protocol

type NewsResp struct {
	Id int64
	Title string
	ThumbUrl string
	Author string
	Avatar string
	ReadCount int
	CommentCount int
	LikeCount int
	Category string
	ViewType int
	IsRecommend int
	Content string
}