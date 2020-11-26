package model

// ResourceRoute 资源Web路由
type ResourceRoute struct {
	Model
	Name             string `json:"name"`
	Path             string `json:"path"`
	Method           string `json:"method"` // 可多选(POST,GET,PUT,DELETE)等
	ResourceServerID ID     `json:"resource_server_id"`
}

// ResourceMenu 资源菜单
type ResourceMenu struct {
	Model
	Name               string          `json:"name" gorm:"column:name"`
	URL                string          `json:"url" gorm:"column:url"`
	Icon               string          `json:"icon" gorm:"column:icon"`
	Level              int             `json:"level" gorm:"column:level"` // 等级、菜单深度
	SerialNumber       int             `json:"serial_number" gorm:"column:serial_number"`
	Leaf               bool            `json:"leaf" gorm:"column:leaf"` // 是：子组件，否：是父组件
	ParentID           ID              `json:"parent_id" gorm:"column:parent_id"`
	ResourceServerID   ID              `json:"resource_server_id" gorm:"column:resource_server_id"`
	ChildResourceMenus []*ResourceMenu `json:"child_resource_menus" gorm:"-"`
}

// ResourceAction web组件
type ResourceAction struct {
	Model
	Code             Code   `json:"code" gorm:"column:code"`
	Name             string `json:"name" gorm:"column:name"`
	Group            string `json:"group" gorm:"column:group"`
	Description      string `json:"description" gorm:"column:description"`
	ResourceServerID ID     `json:"resource_server_id" gorm:"column:resource_server_id"`
}
