package handler

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/harbouriam/configuration"
	redis2 "github.com/harbourrocks/harbour/pkg/harbouriam/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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
}

func DockerScopeFromString(scope string) DockerScope {
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

func DockerScopesFromString(scopes string) []DockerScope {
	split := strings.Split(scopes, " ")
	r := make([]DockerScope, len(split))
	for i, scope := range split {
		r[i] = DockerScopeFromString(scope)
	}
	return r
}

func (h DockerTokenModel) Handle() {
	w := h.GetResponse()
	redisConfig := h.GetRedisConfig()

	qAccount := h.GetQueryParam("account")
	qClientId := h.GetQueryParam("client_id")
	qOfflineToken := strings.ToLower(h.GetQueryParam("offline_token")) == "true"
	qService := h.GetQueryParam("service")
	qScope := h.GetQueryParam("scope")

	dockerScopes := DockerScopesFromString(qScope)

	l.
		WithField("account", qAccount).
		WithField("service", qService).
		WithField("clientId", qClientId).
		WithField("offlineToken", qOfflineToken).
		WithField("scope", qScope).
		Trace("Docker authentication attempt")

	client := redisconfig.OpenClient(redisConfig)

	userNameKey := redis2.IamUserName(qAccount)
	userIdCmd := client.Get(userNameKey)
	if userIdCmd.Err() == redis.Nil {
		// user was not found
		l.WithField("account", qAccount).Warn("Username not found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if userIdCmd.Err() != nil {
		// redis error occurred
		l.WithError(userIdCmd.Err()).WithField("account", qAccount).Error("Failed to load user identification")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := userIdCmd.Val()
	l.WithField("account", qAccount).WithField("userId", userId).Info("Docker authentication attempt")

	dockerPasswordHash, err := client.HGet(redis2.IamUserKey(userId), "docker-password").Result()
	if err != nil {
		// redis error occurred
		l.WithError(userIdCmd.Err()).WithField("userId", userId).Error("Failed to load docker password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.WithField("userId", userId).WithField("hash", dockerPasswordHash).Trace("Docker Password hash found")

	dockerPasswordDecoded, err := base64.StdEncoding.DecodeString(dockerPasswordHash)
	if err != nil {
		l.WithError(err).WithField("userId", userId).Error("Failed to decode docker password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get and validate authorization header value
	var sentUsername, sentPassword string
	authHeaderValue := h.GetRequest().Header.Get("Authorization")
	l.WithField("AuthorizationHeader", authHeaderValue).Trace()
	if len(authHeaderValue) < len("Basic ") {
		l.Warn("AuthorizationHeader value too short")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		authHeaderValue = strings.TrimPrefix(authHeaderValue, "Basic ")
		l.WithField("AuthorizationHeader", authHeaderValue).Trace("Header trimmed")

		bytes, err := base64.StdEncoding.DecodeString(authHeaderValue)
		if err != nil {
			l.WithError(err).Warn("Failed to decode authentication header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authHeaderValue := string(bytes)
		l.WithField("AuthorizationHeader", authHeaderValue).Trace("Header decoded")

		split := strings.Split(authHeaderValue, ":")
		if len(split) != 2 {
			l.Warn("There has to be exactly one colon for basic authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sentUsername, sentPassword = split[0], split[1]
	}

	l.WithField("username", sentUsername).WithField("password", sentPassword).Trace("Extracted username and password")

	// compare passwords
	if bcrypt.CompareHashAndPassword(dockerPasswordDecoded, []byte(sentPassword)) != nil {
		l.WithField("userId", userId).Warn("Invalid password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	l.WithField("userId", userId).Info("Authentication successful")

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

	// read certificate file
	certPath := h.GetIamConfig().Docker.CertificatePath
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		l.WithError(err).WithField("certPath", certPath).Error("Failed to open certificate file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate x5c, we first have to convert from pem to der format, otherwise x509.ParseCertificates throws an error
	block, _ := pem.Decode(certBytes)
	certificates, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		l.WithError(err).WithField("certPath", certPath).Error("Failed to open parse certificate file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate the x5c header which is simple a string array of the certificate chain
	// each certificate is represented by a base64 encoded string
	x5c := make([]string, len(certificates))
	for i, certificate := range certificates {
		x5c[i] = base64.StdEncoding.EncodeToString(certificate.Raw)
	}

	// read private key file
	keyPath := h.GetIamConfig().Docker.SigningKeyPath
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		l.WithError(err).WithField("keyFile", keyPath).Error("Failed to open private key file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// convert key bytes to rsa.PrivateKey
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		l.WithError(err).Error("Failed to parse private key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// add the x5c header
	token.Header["x5c"] = x5c

	// sign the token with the key
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		l.WithError(err).Error("Failed to sign token with private key")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.WithField("token", signedToken).Trace("Signed JWT token")

	w.WriteHeader(http.StatusOK)
	_ = h.WriteResponse(struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		IssuedAt    string `json:"issued_at"`
	}{
		Token:       signedToken,
		AccessToken: signedToken,
		ExpiresIn:   int64(h.GetIamConfig().Docker.TokenLifetime.Seconds()),
		IssuedAt:    nowTime.Format(time.RFC3339),
	})
}
