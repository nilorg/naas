package model

import "encoding/gob"

func init() {
	gob.Register(&SessionAccount{})
	gob.Register(&SessionThirdBind{})
}

type ThirdBindAction string

// SessionAccount ...
type SessionAccount struct {
	UserID   ID
	UserName string
	Nickname string
	Picture  string
}

type SessionThirdBind struct {
	ThirdID string
	Type    UserThirdType
	Extra   interface{}
}
