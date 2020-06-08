package model

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/sdk/convert"
)

// Pagination ...
type Pagination struct {
	Total    uint64 `json:"total"`
	PageSize int    `json:"pageSize"`
	Current  int    `json:"current"`
}

// GetSkip 获取跳过条数
func (p *Pagination) GetSkip() int {
	if p.Current <= 1 {
		return 0
	}
	return (p.Current - 1) * p.PageSize
}

// GetLimit 获取显示条数
func (p *Pagination) GetLimit() int {
	return p.PageSize
}

// NewPagination ...
func NewPagination(ctx *gin.Context) *Pagination {
	return &Pagination{
		Current:  convert.ToInt(ctx.Query("current")),
		PageSize: convert.ToInt(ctx.Query("pageSize")),
		Total:    0,
	}
}

// NewTableListData ...
func NewTableListData(pagination Pagination, list interface{}) *TableListData {
	return &TableListData{
		Pagination: pagination,
		List:       list,
	}
}

// TableListData ...
type TableListData struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}
