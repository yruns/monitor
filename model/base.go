package model

type Page struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}
