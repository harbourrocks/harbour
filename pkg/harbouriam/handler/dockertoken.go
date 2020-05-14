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
	traits.IdTokenModel
}

var supportedScopes = map[string]struct{}{
	"pull": {},
	"push": {},
}

func (h DockerTokenModel) Handle() {
	w := h.GetResponse()
	log := h.GetHRock().L

	qAccount := h.GetQueryParam("account")
	qClientId := h.GetQueryParam("client_id")
	qOfflineToken := strings.ToLower(h.GetQueryParam("offline_token")) == "true"
	qService := h.GetQueryParam("service")
	qScope := h.GetQueryParam("scope")

	log.
		WithField("account", qAccount).
		WithField("service", qService).
		WithField("clientId", qClientId).
		WithField("offlineToken", qOfflineToken).
		WithField("scope", qScope).
		Trace("Docker authentication attempt")

	// authenticate user and resolve dockerUsername
	dockerUserName, err := h.authenticateViaToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if dockerUserName == "" {
		dockerUserName, err = h.authenticateViaBasic()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if dockerUserName == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	log.Info("Authentication successful")

	// validate scopes
	dockerScopes := dockerScopesFromString(qScope)
	for _, scope := range dockerScopes {
		if !validateScopeOk(scope) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	nowTime := time.Now()
	claims := jwt.MapClaims{
		"iss": h.GetIamConfig().Docker.Issuer,                            // this is the issuer (harbour iam namely)
		"sub": dockerUserName,                                            // this is the user who wants to authenticate
		"aud": qService,                                                  // this is the identifier of the registry
		"exp": nowTime.Add(h.GetIamConfig().Docker.TokenLifetime).Unix(), // token expiration
		"nbf": nowTime.Unix(),                                            // not before
		"iat": nowTime.Unix(),                                            // issued at
		"jit": uuid.New().String(),                                       // some unique value required by the registry
	}

	if len(dockerScopes) > 0 {
		claims["access"] = dockerScopes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

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

// dockerScopesFromString converts space separated DockerScope into an array of DockerScopes
func dockerScopesFromString(scopes string) []DockerScope {
	if scopes == "" {
		return make([]DockerScope, 0)
	}

	split := strings.Split(scopes, " ")
	r := make([]DockerScope, len(split))
	for i, scope := range split {
		r[i] = dockerScopeFromString(scope)
	}
	return r
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

func (h DockerTokenModel) getDockerPasswordByUserId(userId string) (dockerPasswordDecoded []byte, err error) {
	log := h.GetHRock().L

	client := redisconfig.OpenClient(h.GetRedisConfig())
	dockerPasswordHash, err := client.HGet(redis2.IamUserKey(userId), "docker-password").Result()
	if err != nil {
		// redis error occurred
		log.WithError(err).Error("Failed to load docker password")
		return
	}

	log.WithField("hash", dockerPasswordHash).Trace("Docker Password hash found")

	dockerPasswordDecoded, err = base64.StdEncoding.DecodeString(dockerPasswordHash)
	if err != nil {
		log.WithError(err).Error("Failed to decode docker password")
		return
	}

	return
}

func (h DockerTokenModel) getDockerUsernameByUserId(userId string) (dockerUsername string, err error) {
	log := h.GetHRock().L

	client := redisconfig.OpenClient(h.GetRedisConfig())
	dockerUsername, err = client.HGet(redis2.IamUserKey(userId), "preferred_username").Result()
	if err != nil {
		// redis error occurred
		log.WithError(err).Error("Failed to load preferred_username")
		return
	}

	return
}

func (h DockerTokenModel) authenticateViaToken() (dockerUsername string, err error) {
	idToken := h.GetHRock().IdToken

	if idToken != nil {
		dockerUsername, err = h.getDockerUsernameByUserId(idToken.Subject)
	}

	return
}

func (h DockerTokenModel) authenticateViaBasic() (dockerUsername string, err error) {
	log := h.GetHRock().L

	authHeaderValue := h.GetRequest().Header.Get("Authorization")
	log.WithField("AuthorizationHeader", authHeaderValue).Trace("Got authorization header")

	// decode basic header
	dockerUsername, sentPassword := decomposeBasicHeader(h.GetHRock(), authHeaderValue)
	log.WithField("username", dockerUsername).WithField("password", sentPassword).Trace("Extracted username and password")
	if dockerUsername == "" {
		return
	}

	// turn username into harbour userId
	userId, err := h.resolveUserIdFromUsername(dockerUsername)
	if err != nil {
		return // error logged in resolveUserIdFromUsername
	} else if userId == "" {
		dockerUsername = ""
		return
	}

	log.WithField("userId", userId).Info("Docker authentication attempt")

	// load docker password hash for user
	dockerPasswordDecoded, err := h.getDockerPasswordByUserId(userId)
	if err != nil {
		return // error logged in getDockerPasswordByUserId
	}

	// compare passwords
	if bcrypt.CompareHashAndPassword(dockerPasswordDecoded, []byte(sentPassword)) != nil {
		log.Warn("Invalid password")
		dockerUsername = ""
		return
	}

	return
}

func decomposeBasicHeader(hRock context.HRock, authHeaderValue string) (username string, password string) {
	log := hRock.L

	// get and validate authorization header value
	if strings.HasPrefix(authHeaderValue, "Basic ") == false {
		log.Warn("No basic authentication header")
		return
	}

	// trim Basic prefix
	authHeaderValue = strings.TrimPrefix(authHeaderValue, "Basic ")
	log.WithField("AuthorizationHeader", authHeaderValue).Trace("Basic header trimmed")

	// decode base64
	bytes, err := base64.StdEncoding.DecodeString(authHeaderValue)
	if err != nil {
		log.WithError(err).Warn("Failed to decode Basic header")
		return
	}

	// bytes to string
	authHeaderValue = string(bytes)
	log.WithField("AuthorizationHeader", authHeaderValue).Trace("Basic header decoded")

	// username and password are separated by a colon
	// if there are two colons then fail
	split := strings.Split(authHeaderValue, ":")
	if len(split) != 2 {
		log.Warn("There has to be exactly one colon for basic authentication")
		return
	}

	username, password = split[0], split[1]
	return
}

type UnauthorizedError struct {
}

func (e UnauthorizedError) Error() string {
	return "authentication failed"
}
