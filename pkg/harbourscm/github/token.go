package github

import (
	"context"
	"errors"
	"github.com/harbourrocks/harbour/pkg/harbourscm/redis"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"strconv"
	"time"
)

// GenerateTokenForOrganization look up a token for the organization
func GenerateTokenForOrganization(ctx context.Context, orgLogin string) (installationToken string, err error) {
	log := logconfig.GetLogCtx(ctx)
	client := redisconfig.GetRedisClientCtx(ctx)

	log.WithField("orgId", orgLogin).Trace("Generating token for organization")

	appIdPemSlice, err := client.HMGet(redis.GithubOrganizationLoginKey(orgLogin), "appId", "pem", "installationId").Result()
	if err != nil {
		log.WithError(err).WithField("orgId", orgLogin).Error("Failed to get appId of organization")
		return
	}

	appIdStr := appIdPemSlice[0].(string)
	privateKeyPem := appIdPemSlice[1].(string)
	installationId := appIdPemSlice[2].(string)

	log.WithField("appId", appIdStr).Trace("Got appId of organization")

	appId, err := strconv.Atoi(appIdStr)
	if appIdStr == "" || err != nil {
		log.WithField("orgLogin", orgLogin).Warn("Missing/invalid appId for organization")
		err = errors.New("missing/invalid appId for organization")
		return
	}

	appToken, err := GenerateGithubAppToken(ctx, appId, privateKeyPem, time.Minute*5)
	if err != nil {
		// error logged in GenerateGithubAppToken
		return
	}

	installationToken, err = GenerateGithubInstallationTokenFromAppToken(ctx, installationId, appToken)
	// error logged in GenerateGithubInstallationTokenFromAppToken

	return
}
