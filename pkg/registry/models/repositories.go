package models

// Repositories is the response models of /v2/_catalog
type Repositories struct {
	Repositories []string `json:"repositories"`
}
