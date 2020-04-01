package model

// User ...
type User struct {
	Model
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}
