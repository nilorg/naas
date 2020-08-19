package model

// Resource 资源
type Resource struct {
	Model
	Name           string `json:"name"`
	Secret         string `json:"secret"`
	Description    string `json:"description"`
	OrganizationID ID     `json:"organization_id"`
}
