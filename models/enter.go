package models

// PageInfo Page 分页
type PageInfo struct {
	Page  int    `form:"page"`
	Key   string `form:"key"`
	Limit int    `form:"limit"`
	Sort  string `form:"sort"`
}

// RemoveRequest 删除
type RemoveRequest struct {
	IDList []uint `json:"id_list"`
}
