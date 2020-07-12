package model

import "time"

// Model ...
type Model struct {
	ID        uint64     `json:"id" gorm:"primary_key;column:id"`
	CreatedAt time.Time  `json:"created_at" sql:"index" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
}

// CodeModel ...
type CodeModel struct {
	Code      string     `json:"code" gorm:"primary_key;column:code"`
	CreatedAt time.Time  `json:"created_at" sql:"index" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
}

// CacheIDPrimaryKey ...
type CacheIDPrimaryKey struct {
	ID uint64 `json:"id"`
}

// CacheCodePrimaryKey ...
type CacheCodePrimaryKey struct {
	Code string `json:"code"`
}
