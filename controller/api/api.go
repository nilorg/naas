package api

var (
	// Admin ..
	Admin = &admin{}
)

// ResultError example...
type ResultError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
