package model

import "github.com/jinzhu/gorm"

// OAuth2Scope ...
type OAuth2Scope struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}
