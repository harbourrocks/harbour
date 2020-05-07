package models

// Repositories is the response model of /v2/_catalog
type Repositories struct {
	Repositories []string `json:"repositories"`
}
