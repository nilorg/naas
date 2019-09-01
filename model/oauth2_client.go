package model

import "github.com/jinzhu/gorm"

// OAuth2Client ...
type OAuth2Client struct {
	gorm.Model
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}
