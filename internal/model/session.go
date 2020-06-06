package model

import "encoding/gob"

func init() {
	gob.Register(&SessionAccount{})
}

// SessionAccount ...
type SessionAccount struct {
	UserID   uint64
	UserName string
	Nickname string
	Picture  string
}
