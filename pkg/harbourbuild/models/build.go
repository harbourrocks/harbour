package models

type Build struct {
	BuildId     string
	Project     string
	Commit      string
	Logs        string
	Repository  string
	BuildStatus string
}
