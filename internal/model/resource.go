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

// ResultResourceRoute 返回资源服务器中的路由
type ResultResourceRoute struct {
	ResourceServer *Resource      `json:"resource_server"`
	ResourceRoute  *ResourceRoute `json:"resource_route"`
}

// ResultResourceMenu 返回资源服务器中的菜单
type ResultResourceMenu struct {
	ResourceServer     *Resource     `json:"resource_server"`
	ResourceMenu       *ResourceMenu `json:"resource_menu"`
	ParentResourceMenu *ResourceMenu `json:"parent_resource_menu"`
}

// ResultResourceAction 返回资源服务器中的动作
type ResultResourceAction struct {
	ResourceServer *Resource       `json:"resource_server"`
	ResourceAction *ResourceAction `json:"resource_action"`
}
