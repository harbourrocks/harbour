package models

import "time"

// GithubAppConfiguration is the app configuration from github
type GithubAppConfiguration struct {
	Id            int         `json:"id"`
	NodeId        string      `json:"node_id"`
	Owner         GithubOwner `json:"owner"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	ExternalUrl   string      `json:"external_url"`
	HtmlUrl       string      `json:"html_url"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	ClientId      string      `json:"client_id"`
	ClientSecret  string      `json:"client_secret"`
	WebhookSecret string      `json:"webhook_secret"`
	PEM           string      `json:"pem"`
}
