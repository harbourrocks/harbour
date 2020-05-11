package models

type BuildRequest struct {
	Dockerfile string   `json:"dockerfile"`
	Tags       []string `json:"tags"`
	Project    string   `json:"project"`
}
