package model

import "encoding/gob"

func init() {
	gob.Register(&SessionAccount{})
	gob.Register(&SessionThirdBind{})
}

type ThirdBindAction string

var (
	SessionAccountActionBindWx ThirdBindAction = "bind_wx"
)

// SessionAccount ...
type SessionAccount struct {
	UserID   ID
	UserName string
	Nickname string
	Picture  string
	WxOpenID string
	Action   ThirdBindAction
}

type SessionThirdBind struct {
	ThirdID string
	Type    UserThirdType
}
