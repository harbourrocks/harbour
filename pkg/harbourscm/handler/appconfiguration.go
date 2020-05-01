package handler

// AppConfiguration is the app configuration from github
type AppConfiguration struct {
	Id            int    `json:"id"`
	NodeId        string `json:"node_id"`
	Name          string `json:"name"`
	ClientId      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	WebhookSecret string `json:"webhook_secret"`
	PEM           string `json:"pem"`
}
