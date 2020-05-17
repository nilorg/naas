package model

import (
	"github.com/nilorg/sdk/convert"
)

// OAuth2Client ...
type OAuth2Client struct {
	ClientID     uint64 `json:"client_id" gorm:"primary_key;column:client_id"`
	ClientSecret string `json:"client_secret" gorm:"column:client_secret"`
	RedirectURI  string `json:"redirect_uri" gorm:"column:redirect_uri"`
}

// TableName ...
func (*OAuth2Client) TableName() string {
	return "oauth2_client"
}

// GetClientID get client id as string.
func (c *OAuth2Client) GetClientID() string {
	return convert.ToString(c.ClientID)
}

// OAuth2ClientInfo ...
type OAuth2ClientInfo struct {
	ClientID    uint64 `gorm:"primary_key;column:client_id"`
	Name        string `json:"name" gorm:"column:name"`
	Profile     string `json:"profile" gorm:"column:profile"`
	Description string `json:"description" gorm:"column:description"`
}

// TableName ...
func (*OAuth2ClientInfo) TableName() string {
	return "oauth2_client_info"
}

// OAuth2Scope ...
type OAuth2Scope struct {
	Model
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

// TableName ...
func (*OAuth2Scope) TableName() string {
	return "oauth2_scope"
}