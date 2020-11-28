package oauth2

import "encoding/json"

// TokenResponse token response.
type TokenResponse struct {
	AccessToken  string      `json:"access_token"`
	TokenType    string      `json:"token_type,omitempty"`
	ExpiresIn    int64       `json:"expires_in"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Scope        string      `json:"scope,omitempty"`
	IDToken      string      `json:"id_token,omitempty"` // https://openid.net/specs/openid-connect-core-1_0.html#IDToken
}

// DeviceAuthorizationResponse Device Authorization Response.
// https://tools.ietf.org/html/rfc8628#section-3.2
type DeviceAuthorizationResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int64  `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// IntrospectionResponse Introspection Response.
// https://tools.ietf.org/html/rfc7662#section-2.2
type IntrospectionResponse struct {
	Active   bool   `json:"active"`
	ClientID string `json:"client_id,omitempty"`
	Username string `json:"username,omitempty"`
	Scope    string `json:"scope,omitempty"`
	Sub      string `json:"sub,omitempty"`
	Aud      string `json:"aud,omitempty"`
	Iss      int64  `json:"iss,omitempty"`
	Exp      int64  `json:"exp,omitempty"`
}

// ErrorResponse error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// CodeValue code值
type CodeValue struct {
	ClientID    string   `json:"client_id"`
	OpenID      string   `json:"open_id"`
	RedirectURI string   `json:"redirect_uri"`
	Scope       []string `json:"scope"`
}

// MarshalBinary json
func (code *CodeValue) MarshalBinary() ([]byte, error) {
	return json.Marshal(code)
}

// UnmarshalBinary json
func (code *CodeValue) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, code)
}

// DeviceCodeValue device_code值
type DeviceCodeValue struct {
	OpenID string   `json:"open_id"`
	Scope  []string `json:"scope"`
}

// MarshalBinary json
func (code *DeviceCodeValue) MarshalBinary() ([]byte, error) {
	return json.Marshal(code)
}

// UnmarshalBinary json
func (code *DeviceCodeValue) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, code)
}

// ClientBasic 客户端基础
type ClientBasic struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
}
