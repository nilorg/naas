package model

// Resource 资源
type Resource struct {
	Model
	Name           string `json:"name"`
	Secret         string `json:"secret"`
	Description    string `json:"description"`
	OrganizationID ID     `json:"organization_id"`
}

// ResultResourceServer 返回资源服务器
type ResultResourceServer struct {
	ResourceServer *Resource     `json:"resource_server"`
	Organization   *Organization `json:"organization"`
}

// ResultResourceWebRoute 返回资源服务器中的Web路由
type ResultResourceWebRoute struct {
	ResourceServer   *Resource         `json:"resource_server"`
	ResourceWebRoute *ResourceWebRoute `json:"resource_web_route"`
}

// ResultResourceWebMenu 返回资源服务器中的Web菜单
type ResultResourceWebMenu struct {
	ResourceServer        *Resource        `json:"resource_server"`
	ResourceWebMenu       *ResourceWebMenu `json:"resource_web_menu"`
	ParentResourceWebMenu *ResourceWebMenu `json:"parent_resource_web_menu"`
}

// ResultResourceAction 返回资源服务器中的动作
type ResultResourceAction struct {
	ResourceServer *Resource       `json:"resource_server"`
	ResourceAction *ResourceAction `json:"resource_action"`
}
