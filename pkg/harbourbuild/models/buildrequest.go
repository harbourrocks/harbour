package models

type BuildRequest struct {
	Dockerfile string `json:"dockerfile"`
	Tag        string `json:"tag"`
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
}
