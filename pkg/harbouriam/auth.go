package harbouriam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/harbourrocks/harbour/pkg/cryptography"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// AuthHandler handles the authentication of an user
type AuthHandler struct {
	OAuthConfig  *oauth2.Config
	OIDCConfig   *oidc.Config
	OIDCProvider *oidc.Provider
	Options      *Options
}

// Redirect redirects a client to the OIDC authentication server
// a state nonce is generated and saved to verify it once the Callback endpoint gets called.
func (a AuthHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	state, err := cryptography.GenerateRandomString(16)

	if err != nil {
		logrus.Error("Failed to generate random string: ", err)
		http.Error(w, "Cryptography error, please try again", http.StatusBadRequest)
		return
	}

	// save state to redis
	redisClient := redisconfig.OpenClient(a.Options.Redis)
	defer redisClient.Close()
	if s := redisClient.Set(fmt.Sprint("state:", state), "", 1*time.Minute); s.Err() != nil {
		logrus.Fatal("Failed to write to redis: ", s.Err())
		http.Error(w, "Database error, please try again", http.StatusBadRequest)
		return
	}

	//param := oauth2.SetAuthURLParam("response_mode", "form_post") // form_post recommended
	redirectURL := a.OAuthConfig.AuthCodeURL(state)
	logrus.Debug("OIDC RedirectURL: ", redirectURL)

	// redirect to authentication server
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// Callback is called by the OIDC authentication server once the user was successfully authenticated
func (a AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// the response from the oidc endpoint must contain a state parameter known by IAM
	// this prevents possible replay attacks because IAM discards
	// a known state parameter the first time it occurs in a response
	serverState := r.URL.Query().Get("state")
	redisClient := redisconfig.OpenClient(a.Options.Redis)

	defer redisClient.Close()
	if s := redisClient.Del(fmt.Sprint("state:", serverState)); s.Val() != 1 {
		logrus.Error("Unknown State: ", s)
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := a.OAuthConfig.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		logrus.Error("Failed to exchange token: ", err)
		http.Error(w, "Failed to exchange token, please try again", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		logrus.Error("No id_token field in oauth2 token.")
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	verifier := a.OIDCProvider.Verifier(a.OIDCConfig)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	logrus.Debugf("Id Token: %s", idToken)
	if err != nil {
		logrus.Error("No id_token field in oauth2 token: ", err)
		http.Error(w, "Failed to verify ID Token, please try again", http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalIndent(oauth2Token, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
