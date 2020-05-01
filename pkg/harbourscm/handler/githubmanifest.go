package handler

type GithubManifest struct {
	Name               string             `json:"name"`
	Url                string             `json:"url"`
	HookAttributes     HookAttributes     `json:"hook_attributes"`
	RedirectUrl        string             `json:"redirect_url"`
	DefaultPermissions DefaultPermissions `json:"default_permissions"`
	DefaultEvents      []GithubEventType  `json:"default_events"`
}

type HookAttributes struct {
	Url    string `json:"url"`
	Active bool   `json:"active"`
}

type DefaultPermissions struct {
	Contents GithubPermissionType `json:"contents"`
	Metadata GithubPermissionType `json:"metadata"`
}

type GithubEventType string

const (
	Push GithubEventType = "push"
)

type GithubPermissionType string

const (
	NoAccess GithubPermissionType = "no access"
	Read     GithubPermissionType = "read-only"
	Write    GithubPermissionType = "write"
)
