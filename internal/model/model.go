package model

import (
	"time"

	"github.com/nilorg/naas/internal/pkg/diff"
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

// ConvertIDSliceToInt64Slice id slice 转 int64 slice
func ConvertIDSliceToInt64Slice(ids []ID) (v []int64) {
	for _, id := range ids {
		v = append(v, int64(id))
	}
	return v
}

// ConvertInt64SliceToIDSlice int64 slice 转 id slice
func ConvertInt64SliceToIDSlice(ints []int64) (v []ID) {
	for _, id := range ints {
		v = append(v, ID(id))
	}
	return v
}

// ConvertStringToID string 转 ID
func ConvertStringToID(id string) ID {
	return ID(convert.ToUint64(id))
}

// ConvertStringToCode string 转 code
func ConvertStringToCode(code string) Code {
	return Code(code)
}

// DiffIDSlice ...
func DiffIDSlice(src []ID, in []ID) (added []ID, deleted []ID) {
	srcInts := ConvertIDSliceToInt64Slice(src)
	inInts := ConvertIDSliceToInt64Slice(in)
	addedInts, deletedInts := diff.IntSlice(srcInts, inInts)
	added = ConvertInt64SliceToIDSlice(addedInts)
	deleted = ConvertInt64SliceToIDSlice(deletedInts)
	return
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
