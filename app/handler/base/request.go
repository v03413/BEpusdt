package base

// ListRequest 通用列表请求
type ListRequest struct {
	Page    int    `json:"page" form:"page" binding:"min=1"`
	Size    int    `json:"size" form:"size" binding:"min=1"`
	Keyword string `json:"keyword" form:"keyword"`
	Sort    string `json:"sort" form:"sort" binding:"oneof=asc desc"`
	ID      int    `json:"id" form:"id"`
	Status  int    `json:"status" form:"status"`
}

type IDRequest struct {
	ID int `json:"id" form:"id" binding:"required"` // *ID(必须)
}
