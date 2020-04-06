package harbouriam

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/coreos/go-oidc"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// RunIAMServer runns the IAM server application
func RunIAMServer(o *Options) error {
	logrus.Info("Started Harbour IAM server")

	ctx := context.Background()

	// discover oidc endpoint configuration
	provider, err := oidc.NewProvider(ctx, o.OIDCURL)
	if err != nil {
		logrus.Fatal(err)
	}

	oidcConfig := &oidc.Config{
		ClientID: o.OIDCClientID,
	}

	// obtain login redirect url
	redirectURL, err := url.Parse(o.IAMBaseURL)
	if err != nil {
		logrus.Fatal(err)
	} else {
		redirectURL.Path = path.Join(redirectURL.Path, "/auth/oidc/callback")
	}

	config := &oauth2.Config{
		ClientID:     o.OIDCClientID,
		ClientSecret: o.OIDCClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL.String(),
		Scopes:       []string{oidc.ScopeOpenID, "email"},
	}

	authHandler := AuthHandler{
		OAuthConfig:  config,
		OIDCConfig:   oidcConfig,
		OIDCProvider: provider,
		Options:      o,
	}

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
		authHandler.Redirect(w, r)
	})

	http.HandleFunc("/auth/oidc/callback", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
		authHandler.Callback(w, r)
	})

	bindAddress := "127.0.0.1:5100"
	logrus.Info(fmt.Sprintf("Listening on http://%s/", bindAddress))

	err = http.ListenAndServe(bindAddress, nil)
	logrus.Fatal(err)

	return err
}
