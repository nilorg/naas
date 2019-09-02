package model

import (
	"encoding/gob"
)

func init() {
	gob.Register(&User{})
}

// User ...
type User struct {
	Model
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}
