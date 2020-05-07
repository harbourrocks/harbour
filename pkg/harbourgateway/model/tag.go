package model

// Tag represents a tagged image in a specific repository
//  A tag is identified by its name
type Tag struct {
	Name       string     `json:"Name"`
	Repository Repository `json:"repository"`
}
