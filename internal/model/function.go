package model

const (
	// WebFunctionTypeMenu ...
	WebFunctionTypeMenu = ""
)

// WebFunction web功能
type WebFunction struct {
	Model
	Name        string         `json:"name" gorm:"column:name"`
	Type        string         `json:"type" gorm:"column:type"`
	Description string         `json:"description" gorm:"column:description"`
	Components  []WebComponent `json:"components" gorm:"-"`
	Extension   string         `json:"extension" gorm:"column:extension"`
}

// WebFunctionComponent web功能中的组件
type WebFunctionComponent struct {
	Model
	FunctionID  uint64 `json:"function_id" gorm:"column:function_id"`
	ComponentID uint64 `json:"component_id" gorm:"column:component_id"`
}
