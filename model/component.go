package model

const (
	// WebComponentTypeInput ...
	WebComponentTypeInput = "input"
	// WebComponentTypeButton ...
	WebComponentTypeButton = "button"
)

// WebComponent web组件
type WebComponent struct {
	Model
	Name        string `json:"name" gorm:"column:name"`
	Type        string `json:"type" gorm:"column:type"`
	TypeValue   string `json:"type_value" gorm:"column:type_value"`
	Description string `json:"description" gorm:"column:description"`
	Extension   string `json:"extension" gorm:"column:extension"`
}
