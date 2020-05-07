package model

// Repository represents a repository on the docker registry
//  A repository is identified by its name
type Repository struct {
	Name string `json:"name"`
}
