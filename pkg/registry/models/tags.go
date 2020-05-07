package models

// Tags is the response model of /v2/<repository-name>/tags/list
type Tags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
