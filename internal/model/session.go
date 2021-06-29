package model

import "encoding/gob"

func init() {
	gob.Register(&SessionAccount{})
}

type SessionAccountAction string

var (
	SessionAccountActionBindWx SessionAccountAction = "bind_wx"
)

// SessionAccount ...
type SessionAccount struct {
	UserID   ID
	WxOpenID string
	UserName string
	Nickname string
	Picture  string
	Action   SessionAccountAction
}
