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
	UserID uint64 `gorm:"not null"`
	// 昵称
	NickName string `gorm:"null"`
	// 头像
	AvatarURL string `gorm:"null"`
	// 性别
	Gender int `gorm:"null"`
}
