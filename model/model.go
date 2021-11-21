package model

import (
	"time"
)

// Model Customized gorm.Model
type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// PaginationQuery 更新信息查询
type PaginationQuery struct {
	// Page 第ｎ页
	Page int `json:"page,omitempty" form:"page"`
	// PageSize 每页的大小
	PageSize int `json:"page_size,omitempty" form:"page_size"`
	// Sort 排序字段
	Sort string `json:"sort,omitempty" form:"sort"`
}
