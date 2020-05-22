package models

import "time"

type GithubRepository struct {
	Id                 int               `json:"id"`
	NodeId             string            `json:"node_id"`
	Name               string            `json:"name"`
	FullName           string            `json:"full_name"`
	Owner              GithubOwner       `json:"owner"`
	Private            bool              `json:"private"`
	HTMLUrl            string            `json:"html_url"`
	Description        string            `json:"description"`
	Fork               bool              `json:"fork"`
	Url                string            `json:"url"`
	ArchiveUrl         string            `json:"archive_url"`
	AssigneesUrl       string            `json:"assignees_url"`
	BlobsUrl           string            `json:"blobs_url"`
	BranchesUrl        string            `json:"branches_url"`
	CollaboratorsUrl   string            `json:"collaborators_url"`
	CommentsUrl        string            `json:"comments_url"`
	CommitsUrl         string            `json:"commits_url"`
	CompareUrl         string            `json:"compare_url"`
	ContentsUrl        string            `json:"contents_url"`
	ContributorsUrl    string            `json:"contributors_url"`
	DeploymentsUrl     string            `json:"deployments_url"`
	DownloadsUrl       string            `json:"downloads_url"`
	EventsUrl          string            `json:"events_url"`
	ForksUrl           string            `json:"forks_url"`
	GitCommitsUrl      string            `json:"git_commits_url"`
	GitRefsUrl         string            `json:"git_refs_url"`
	GitTagsUrl         string            `json:"git_tags_url"`
	GitUrl             string            `json:"git_url"`
	IssueCommentUrl    string            `json:"issue_comment_url"`
	IssueEventsUrl     string            `json:"issue_events_url"`
	IssuesUrl          string            `json:"issues_url"`
	KeysUrl            string            `json:"keys_url"`
	LabelsUrl          string            `json:"labels_url"`
	LanguagesUrl       string            `json:"languages_url"`
	MergesUrl          string            `json:"merges_url"`
	MilestonesUrl      string            `json:"milestones_url"`
	NotificationsUrl   string            `json:"notifications_url"`
	PullsUrl           string            `json:"pulls_url"`
	ReleasesUrl        string            `json:"releases_url"`
	SSHUrl             string            `json:"ssh_url"`
	StargazersUrl      string            `json:"stargazers_url"`
	StatusesUrl        string            `json:"statuses_url"`
	SubscribersUrl     string            `json:"subscribers_url"`
	SubscriptionUrl    string            `json:"subscription_url"`
	TagsUrl            string            `json:"tags_url"`
	TeamsUrl           string            `json:"teams_url"`
	TreesUrl           string            `json:"trees_url"`
	CloneUrl           string            `json:"clone_url"`
	MirrorUrl          string            `json:"mirror_url"`
	HooksUrl           string            `json:"hooks_url"`
	SvnUrl             string            `json:"svn_url"`
	Homepage           string            `json:"homepage"`
	Language           interface{}       `json:"language"`
	ForksCount         int               `json:"forks_count"`
	StargazersCount    int               `json:"stargazers_count"`
	WatchersCount      int               `json:"watchers_count"`
	Size               int               `json:"size"`
	DefaultBranch      string            `json:"default_branch"`
	OpenIssuesCount    int               `json:"open_issues_count"`
	IsTemplate         bool              `json:"is_template"`
	Topics             []string          `json:"topics"`
	HasIssues          bool              `json:"has_issues"`
	HasProjects        bool              `json:"has_projects"`
	HasWiki            bool              `json:"has_wiki"`
	HasPages           bool              `json:"has_pages"`
	HasDownloads       bool              `json:"has_downloads"`
	Archived           bool              `json:"archived"`
	Disabled           bool              `json:"disabled"`
	Visibility         string            `json:"visibility"`
	PushedAt           time.Time         `json:"pushed_at"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	Permissions        GithubPermissions `json:"permissions"`
	TemplateRepository interface{}       `json:"template_repository"`
	TempCloneToken     string            `json:"temp_clone_token"`
	SubscribersCount   int               `json:"subscribers_count"`
	NetworkCount       int               `json:"network_count"`
	License            GithubLicense     `json:"license"`
}

type GithubPermissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type GithubLicense struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxId string `json:"spdx_id"`
	Url    string `json:"url"`
	NodeId string `json:"node_id"`
}
