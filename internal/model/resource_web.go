package model

// ResourceWebRoute 资源Web路由
type ResourceWebRoute struct {
	Model
	Name       string `json:"name"`
	Path       string `json:"path"`
	Method     string `json:"method"` // 可多选(POST,GET,PUT,DELETE)等
	ResourceID ID     `json:"resource_id"`
}

const (
	// WebComponentTypeInput ...
	WebComponentTypeInput = "input"
	// WebComponentTypeButton ...
	WebComponentTypeButton = "button"
)

// ResourceWebComponent web组件
type ResourceWebComponent struct {
	Model
	Name        string `json:"name" gorm:"column:name"`
	Type        string `json:"type" gorm:"column:type"`
	TypeValue   string `json:"type_value" gorm:"column:type_value"`
	Description string `json:"description" gorm:"column:description"`
	Extension   string `json:"extension" gorm:"column:extension"`
}

const (
	// WebFunctionTypeMenu ...
	WebFunctionTypeMenu = "menu"
)

// ResourceWebFunction web功能
type ResourceWebFunction struct {
	Model
	Name        string                 `json:"name" gorm:"column:name"`
	Type        string                 `json:"type" gorm:"column:type"`
	Description string                 `json:"description" gorm:"column:description"`
	Components  []ResourceWebComponent `json:"components" gorm:"-"`
	Extension   string                 `json:"extension" gorm:"column:extension"`
}

// ResourceWebFunctionComponent web功能中的组件
type ResourceWebFunctionComponent struct {
	Model
	FunctionID  ID `json:"function_id" gorm:"column:function_id"`
	ComponentID ID `json:"component_id" gorm:"column:component_id"`
}
