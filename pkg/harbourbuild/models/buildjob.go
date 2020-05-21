package models

// BuildJob represents a job send to the Builder with a channel.
// Contains all information required for build a image.
type BuildJob struct {
	Request       BuildRequest
	BuildKey      string
	RegistryToken string
	RegistryUrl   string
	ReqId         string
}
