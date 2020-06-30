package model

// Resource 资源
type Resource struct {
	Model
	Name           string `json:"name"`
	Secret         string `json:"secret"`
	Description    string `json:"description"`
	OrganizationID uint64 `json:"organization_id"`
}
