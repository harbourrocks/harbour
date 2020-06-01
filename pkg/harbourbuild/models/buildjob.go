package models

// BuildJob represents a job send to the Builder with a channel.
// Contains all information required for build a image.
type BuildJob struct {
	Repository    string
	Tag           string
	FilePath      string
	Dockerfile    string
	BuildKey      string
	RegistryToken string
	RegistryUrl   string
	ReqId         string
}
