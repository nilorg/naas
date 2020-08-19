package model

import (
	"time"

	"github.com/nilorg/sdk/convert"
)

type (
	// ID ...
	ID uint64
	// Code ...
	Code string
)

// ConvertCodeSliceToStringSlice code slice 转 string slice
func ConvertCodeSliceToStringSlice(codes []Code) (v []string) {
	for _, code := range codes {
		v = append(v, string(code))
	}
	return v
}

// ConvertStringToID string 转 ID
func ConvertStringToID(id string) ID {
	return ID(convert.ToUint64(id))
}

// Model ...
type Model struct {
	ID        ID         `json:"id" gorm:"primary_key;column:id"`
	CreatedAt time.Time  `json:"created_at" sql:"index" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
}

// CodeModel ...
type CodeModel struct {
	Code      Code       `json:"code" gorm:"primary_key;column:code"`
	CreatedAt time.Time  `json:"created_at" sql:"index" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
}

// CacheIDPrimaryKey ...
type CacheIDPrimaryKey struct {
	ID ID `json:"id"`
}

// CacheCodePrimaryKey ...
type CacheCodePrimaryKey struct {
	Code Code `json:"code"`
}
