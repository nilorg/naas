package model

import "time"

// Model ...
type Model struct {
	ID        uint64     `gorm:"primary_key;column:id"`
	CreatedAt time.Time  `sql:"index" gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `sql:"index" gorm:"column:deleted_at"`
}
