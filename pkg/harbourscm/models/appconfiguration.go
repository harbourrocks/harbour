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

type GithubOwner struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLUrl           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
