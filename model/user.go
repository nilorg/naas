package model

import (
	"encoding/gob"

	"github.com/jinzhu/gorm"
)

func init() {
	gob.Register(&User{})
}

// User ...
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}
