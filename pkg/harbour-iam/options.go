package harbouriam

// Options defines all options available to configure the IAM server.
type Options struct {
	OIDCClientID      string
	OIDCClientSecret  string
	OIDCUrl           string
	OIDCLoginCallback string
}

// NewDefaultOptions returns the default options
func NewDefaultOptions() *Options {
	s := Options{
		OIDCClientID:      "",
		OIDCClientSecret:  "",
		OIDCUrl:           "",
		OIDCLoginCallback: "",
	}

	return &s
}
