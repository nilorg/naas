package model

type UserThirdType string

var (
	UserThirdTypePhone           UserThirdType = "phone"
	UserThirdTypeWxUnionID       UserThirdType = "wx_union_id"
	UserThirdTypeWxOpenIDForKfpt UserThirdType = "wxkfpt_open_id"
	UserThirdTypeWxOpenIDForFwh  UserThirdType = "wxfwh_open_id"
	UserThirdTypeWxOpenIDForDyh  UserThirdType = "wxdyh_open_id"
)

// 用户第三方关联
type UserThird struct {
	Model
	UserID     ID            `json:"user_id" gorm:"column:user_id"`         // 用户ID
	ThirdType  UserThirdType `json:"third_type" gorm:"column:third_type"`   // 第三方类型
	ThirdID    string        `json:"third_id" gorm:"column:third_id"`       // 第三方关联ID
	ThirdExtra string        `json:"third_extra" gorm:"column:third_extra"` // 第三方扩展
}

// UserThirdTypeVerify 用户第三方类型验证
func UserThirdTypeVerify(typ UserThirdType) bool {
	switch typ {
	case UserThirdTypeWxUnionID, UserThirdTypeWxOpenIDForKfpt:
		return true
	default:
		return false
	}
}
