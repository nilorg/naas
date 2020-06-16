package model

// Scope 范围
type Scope struct {
	CodeModel
	Name        string `json:"name"`
	Description string `json:"description"`
}

// OAuth2ClientScope OAuth2客户端范围
type OAuth2ClientScope struct {
	Model
	OAuth2ClientID uint64 `json:"oauth2_client_id" gorm:"column:oauth2_client_id"`
	ScopeCode      string `json:"scope_code" gorm:"column:scope_code"`
}
