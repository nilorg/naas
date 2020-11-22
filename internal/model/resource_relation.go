package model

// ResourceWebRoute 资源Web路由
type ResourceWebRoute struct {
	Model
	Name             string `json:"name"`
	Path             string `json:"path"`
	Method           string `json:"method"` // 可多选(POST,GET,PUT,DELETE)等
	ResourceServerID ID     `json:"resource_server_id"`
}

// ResourceWebMenu 资源Web菜单
type ResourceWebMenu struct {
	Model
	Name                  string             `json:"name" gorm:"column:name"`
	URL                   string             `json:"url" gorm:"column:url"`
	Icon                  string             `json:"icon" gorm:"column:icon"`
	Level                 int                `json:"level" gorm:"column:level"` // 等级、菜单深度
	SerialNumber          int                `json:"serial_number" gorm:"column:serial_number"`
	Leaf                  bool               `json:"leaf" gorm:"column:leaf"` // 是：子组件，否：是父组件
	ParentID              ID                 `json:"parent_id" gorm:"column:parent_id"`
	ResourceServerID      ID                 `json:"resource_server_id" gorm:"column:resource_server_id"`
	ChildResourceWebMenus []*ResourceWebMenu `json:"child_resource_web_menus" gorm:"-"`
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
