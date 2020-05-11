package handler

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/context"
	"github.com/harbourrocks/harbour/pkg/cryptography"
	"github.com/harbourrocks/harbour/pkg/harbouriam/configuration"
	redis2 "github.com/harbourrocks/harbour/pkg/harbouriam/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

// DockerTokenResponse is the response model of the token request
type DockerTokenResponse struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	IssuedAt    string `json:"issued_at"`
}

// DockerScope is part of the claims and represents allowed actions
type DockerScope struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Actions []string `json:"actions"`
}

// DockerTokenModel is specific for one handler
type DockerTokenModel struct {
	traits.HttpModel
	redisconfig.RedisModel
	configuration.IamModel
	context.HRockModel
}

var supportedScopes = map[string]struct{}{
	"pull": {},
	"push": {},
}

// dockerScopeFromString converts a scope string into a DockerScope
func dockerScopeFromString(scope string) DockerScope {
	split := strings.Split(scope, ":")
	resourceType := split[0]
	var resourceName, actions string

	if len(split) == 3 {
		resourceName = split[1]
		actions = split[2]
	} else {
		resourceName = split[1] + split[2]
		actions = split[3]
	}

	return DockerScope{
		Type:    resourceType,
		Name:    resourceName,
		Actions: strings.Split(actions, ","),
	}
}

// dockerScopesFromString converts space separated DockerScope into an array of DockerScopes
func dockerScopesFromString(scopes string) []DockerScope {
	split := strings.Split(scopes, " ")
	r := make([]DockerScope, len(split))
	for i, scope := range split {
		r[i] = dockerScopeFromString(scope)
	}
	return r
}

// validateScopeOk checks if the requested scope is okay
func validateScopeOk(scope DockerScope) bool {
	for _, action := range scope.Actions {
		if _, ok := supportedScopes[action]; !ok {
			l.WithField("Action", action).Warn("Action not supported")
			return false
		}
	}

	return true
}

// resolveUserIdFromUsername returns the userId (harbour userId) from a username (via redis lookup)
func (h DockerTokenModel) resolveUserIdFromUsername(username string) (userId string, err error) {
	log := h.GetHRock().L
	client := redisconfig.OpenClient(h.GetRedisConfig())

	userIdCmd := client.Get(redis2.IamUserName(username))
	if userIdCmd.Err() == redis.Nil {
		log.WithField("username", username).Warn("Username not found")
		return
	} else if userIdCmd.Err() != nil {
		log.WithError(userIdCmd.Err()).Error("Failed to load user identification")
		return
	} else {
		userId = userIdCmd.Val()
	}

	return
}

func (h DockerTokenModel) Handle() {
	w := h.GetResponse()
	redisConfig := h.GetRedisConfig()
	log := h.GetHRock().L

	qAccount := h.GetQueryParam("account")
	qClientId := h.GetQueryParam("client_id")
	qOfflineToken := strings.ToLower(h.GetQueryParam("offline_token")) == "true"
	qService := h.GetQueryParam("service")
	qScope := h.GetQueryParam("scope")

	dockerScopes := dockerScopesFromString(qScope)

	// validate scopes
	for _, scope := range dockerScopes {
		if !validateScopeOk(scope) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	log.
		WithField("account", qAccount).
		WithField("service", qService).
		WithField("clientId", qClientId).
		WithField("offlineToken", qOfflineToken).
		WithField("scope", qScope).
		Trace("Docker authentication attempt")

	client := redisconfig.OpenClient(redisConfig)

	userId, err := h.resolveUserIdFromUsername(qAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // error logged in resolveUserIdFromUsername
	} else if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.WithField("userId", userId).Info("Docker authentication attempt")

	dockerPasswordHash, err := client.HGet(redis2.IamUserKey(userId), "docker-password").Result()
	if err != nil {
		// redis error occurred
		log.WithError(err).Error("Failed to load docker password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithField("hash", dockerPasswordHash).Trace("Docker Password hash found")

	dockerPasswordDecoded, err := base64.StdEncoding.DecodeString(dockerPasswordHash)
	if err != nil {
		log.WithError(err).Error("Failed to decode docker password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get and validate authorization header value
	var sentUsername, sentPassword string
	authHeaderValue := h.GetRequest().Header.Get("Authorization")
	log.WithField("AuthorizationHeader", authHeaderValue).Trace("Got authorization header")
	if len(authHeaderValue) < len("Basic ") {
		log.Warn("AuthorizationHeader value too short")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		authHeaderValue = strings.TrimPrefix(authHeaderValue, "Basic ")
		log.WithField("AuthorizationHeader", authHeaderValue).Trace("Authorization header trimmed")

		bytes, err := base64.StdEncoding.DecodeString(authHeaderValue)
		if err != nil {
			log.WithError(err).Warn("Failed to decode authentication header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authHeaderValue := string(bytes)
		log.WithField("AuthorizationHeader", authHeaderValue).Trace("Authorization header decoded")

		split := strings.Split(authHeaderValue, ":")
		if len(split) != 2 {
			log.Warn("There has to be exactly one colon for basic authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sentUsername, sentPassword = split[0], split[1]
	}

	log.WithField("username", sentUsername).WithField("password", sentPassword).Trace("Extracted username and password")

	// compare passwords
	if bcrypt.CompareHashAndPassword(dockerPasswordDecoded, []byte(sentPassword)) != nil {
		log.Warn("Invalid password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Info("Authentication successful")

	nowTime := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":    h.GetIamConfig().Docker.Issuer,                            // this is the issuer (harbour iam namely)
		"sub":    sentUsername,                                              // this is the user who wants to authenticate
		"aud":    qService,                                                  // this is the identifier of the registry
		"exp":    nowTime.Add(h.GetIamConfig().Docker.TokenLifetime).Unix(), // token expiration
		"nbf":    nowTime.Unix(),                                            // not before
		"iat":    nowTime.Unix(),                                            // issued at
		"jit":    uuid.New().String(),                                       // some unique value required by the registry
		"access": dockerScopes,
	})

	// generate the x5c signature
	x5c, err := cryptography.GenerateX5C(h.GetHRock(), h.GetIamConfig().Docker.CertificatePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // error logged in GenerateX5C
	}

	// add the x5c header
	token.Header["x5c"] = x5c

	// read private key file
	privateKey, err := cryptography.ReadPrivateKey(h.GetHRock(), h.GetIamConfig().Docker.SigningKeyPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // error logged in ReadPrivateKey
	}

	// sign the token with the key
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.WithError(err).Error("Failed to sign token with private key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithField("token", signedToken).Trace("Signed JWT token")

	w.WriteHeader(http.StatusOK)
	_ = h.WriteResponse(DockerTokenResponse{
		Token:       signedToken,
		AccessToken: signedToken,
		ExpiresIn:   int64(h.GetIamConfig().Docker.TokenLifetime.Seconds()),
		IssuedAt:    nowTime.Format(time.RFC3339),
	})
}
