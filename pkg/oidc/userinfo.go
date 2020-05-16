package oidc

// Userinfo 用户信息
// https://openid.net/specs/openid-connect-core-1_0.html#Claims
type Userinfo struct {
	Sub                 string           `json:"sub"`                             // 权限范围 openid
	Name                string           `json:"name,omitempty"`                  // 权限范围 profile
	FamilyName          string           `json:"family_name,omitempty"`           // 权限范围 profile
	GivenName           string           `json:"given_name,omitempty"`            // 权限范围 profile
	MiddleName          string           `json:"middle_name,omitempty"`           // 权限范围 profile
	Nickname            string           `json:"nickname,omitempty"`              // 权限范围 profile
	PreferredUsername   string           `json:"preferred_username,omitempty"`    // 权限范围 profile
	Profile             string           `json:"profile,omitempty"`               // 权限范围 profile
	Picture             string           `json:"picture,omitempty"`               // 权限范围 profile
	Website             string           `json:"website,omitempty"`               // 权限范围 profile
	Gender              string           `json:"gender,omitempty"`                // 权限范围 profile
	Birthdate           string           `json:"birthdate,omitempty"`             // 权限范围 profile
	Zoneinfo            string           `json:"zoneinfo,omitempty"`              // 权限范围 profile
	Locale              string           `json:"locale,omitempty"`                // 权限范围 profile
	UpdatedAt           int64            `json:"updated_at,omitempty"`            // 权限范围 profile
	Email               string           `json:"email,omitempty"`                 // 权限范围 email
	EmailVerified       bool             `json:"email_verified,omitempty"`        // 权限范围 email
	PhoneNumber         string           `json:"phone_number,omitempty"`          // 权限范围 phone
	PhoneNumberVerified string           `json:"phone_number_verified,omitempty"` // 权限范围 phone
	Address             *UserinfoAddress `json:"address,omitempty"`               // 权限范围 address
}

// UserinfoAddress 用户地址
// 权限范围 address
type UserinfoAddress struct {
	Formatted     string `json:"formatted,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Region        string `json:"region,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Country       string `json:"country,omitempty"`
}
