package model

type Page struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}

type PageResult struct {
	Data  interface{} `form:"data" json:"data"`
	Total int64       `form:"total" json:"total"`
}
