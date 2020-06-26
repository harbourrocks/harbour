package model

type DockerSetPasswordRequest struct {
	Password string `json:"password"`
}

type DockerSetPasswordResponse struct {
	Username    string `json:"username"`
	PasswordSet bool   `json:"passwordSet"`
}
