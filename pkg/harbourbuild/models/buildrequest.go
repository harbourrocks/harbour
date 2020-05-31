package models

type BuildRequest struct {
	Dockerfile string `json:"dockerfile"`
	Tag        string `json:"tag"`
	SCMId      string `json:"scm_id"`
	Commit     string `json:"commit"`
}
