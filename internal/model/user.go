package model

// User ...
type User struct {
	Model
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

//UserInfo 用户信息
type UserInfo struct {
	Model
	UserID    uint64 `json:"user_id" gorm:"not null"`
	NickName  string `json:"nick_name" gorm:"null"`
	AvatarURL string `json:"avatar_url" gorm:"null"`
	Gender    int    `json:"gender" gorm:"null"`
}
