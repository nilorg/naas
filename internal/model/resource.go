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
