package model

// Resource 资源
type Resource struct {
	Model
	Name           string `json:"name"`
	Secret         string `json:"secret"`
	Description    string `json:"description"`
	OrganizationID ID     `json:"organization_id"`
}

// ResultResource 返回资源
type ResultResource struct {
	Resource     *Resource     `json:"resource"`
	Organization *Organization `json:"organization"`
}
