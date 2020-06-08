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
	UserID        uint64 `json:"user_id" gorm:"column:user_id"`
	Nickname      string `json:"nickname" gorm:"column:nickname"`
	Picture       string `json:"picture" gorm:"column:picture;type:MEDIUMTEXT"` // 支持存储base64头像
	Gender        int    `json:"gender" gorm:"column:gender"`
	Email         string `json:"email" gorm:"column:email"`
	EmailVerified bool   `json:"email_verified" gorm:"column:email_verified"`
	Phone         string `json:"phone" gorm:"column:phone"`
	PhoneVerified bool   `json:"phone_verified" gorm:"column:phone_verified"`
}

// ResultUserInfo 返回用户信息
type ResultUserInfo struct {
	*User
	*UserInfo
}
