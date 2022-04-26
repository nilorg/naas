package model

import "encoding/gob"

func init() {
	gob.Register(&SessionAccount{})
	gob.Register(&SessionThirdBind{})
}

type ThirdBindAction string

// SessionAccount ...
type SessionAccount struct {
	UserID   ID     `json:"user_id"`
	UserName string `json:"user_name"`
	Nickname string `json:"nick_name"`
	Picture  string `json:"picture"`
}

type SessionThirdBind struct {
	ThirdID string
	Type    UserThirdType
	Extra   interface{}
}
