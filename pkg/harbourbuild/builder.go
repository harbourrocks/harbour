package harbourbuild

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/jhoonb/archivex"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Builder struct {
	jobChan      chan models.BuildJob
	cli          *client.Client
	ctx          context.Context
	ctxPath      string
	repoPath     string
	redisOptions redisconfig.RedisOptions
}

func NewBuilder(jobChan chan models.BuildJob, ctxPath string, repoPath string, redisConfig redisconfig.RedisOptions) (Builder, error) {
	var builder Builder
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return builder, err
	}

	builder = Builder{jobChan: jobChan, cli: cli, ctx: ctx, ctxPath: ctxPath, repoPath: repoPath, redisOptions: redisConfig}
	return builder, nil
}

func (b Builder) Start() {
	go func() {
		for {
			select {
			case job := <-b.jobChan:
				b.buildImage(job)
			}
		}
	}()
}

func (b Builder) buildImage(job models.BuildJob) {
	log := logrus.WithField("reqId", job.ReqId)
	redisClient := redisconfig.OpenClient(b.redisOptions)

	if err := redisClient.HSet(job.BuildKey, "build_status", "Running").Err(); err != nil {
		log.WithError(err).Error("Failed to save data to redis")
		return
	}

	buildCtx, err := b.createBuildContext(job.Request.Repository)
	if err != nil {
		log.WithError(err).Error("Failed to create build context")
		return
	}

	tag := []string{b.getImageString(job.RegistryUrl, job.Request.Repository, job.Request.Tag)}
	opt := types.ImageBuildOptions{
		Tags:       tag,
		Dockerfile: job.Request.Dockerfile,
	}

	resp, err := b.cli.ImageBuild(b.ctx, buildCtx, opt)
	if err != nil {
		log.WithError(err).Error("Failed to build image")
		return
	}

	defer func() {
		err := buildCtx.Close()
		err = os.Remove(buildCtx.Name())
		err = resp.Body.Close()
		if err != nil {
			log.WithError(err).Error("Error while cleaning up build context")
			return
		}
	}()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to parse response body")
	}

	logs := buf.String()
	if strings.Contains(logs, "errorDetail") {
		if err := redisClient.HSet(job.BuildKey, "build_status", "Failed", "logs", logs).Err(); err != nil {
			log.WithError(err).Error("Failed to save data to redis")
			return
		}
		return
	}

	if err := redisClient.HSet(job.BuildKey, "build_status", "Success", "logs", logs).Err(); err != nil {
		log.WithError(err).Error("Failed to save data to redis")
		return
	}
	log.Tracef("Image %s was built", job.Request.Repository)

	imageString := b.getImageString(job.RegistryUrl, job.Request.Repository, job.Request.Tag)

	if err = b.pushImage(imageString, job.RegistryToken); err != nil {
		log.WithError(err).Error("Error while pushing image to registry")
		return
	}

	log.Tracef("Image %s was pushed to registry %s", job.Request.Repository, job.RegistryUrl)
}

func (b Builder) createBuildContext(project string) (*os.File, error) {
	buildContext := fmt.Sprintf(b.ctxPath+"%s.tar", project)
	projectPath, err := b.getProjectPath(project)
	if err != nil {
		return nil, err
	}

	tar := new(archivex.TarFile)
	err = tar.Create(buildContext)
	err = tar.AddAll(projectPath, false)
	err = tar.Close()

	dockerBuildCtx, err := os.Open(buildContext)
	if err != nil {
		return nil, err
	}

	return dockerBuildCtx, nil
}

func (b Builder) pushImage(image string, token string) error {
	authConfig := types.AuthConfig{RegistryToken: token}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	options := types.ImagePushOptions{RegistryAuth: authStr}

	out, err := b.cli.ImagePush(b.ctx, image, options)
	if err != nil {
		return err
	}

	defer func() {
		err = out.Close()
		if err != nil {
			return
		}
	}()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(out)
	if err != nil {
		return err
	}

	logs := buf.String()
	if strings.Contains(logs, "errorDetail") {
		return fmt.Errorf("unable to push image: %s", logs)
	}

	return nil
}

func (b Builder) getImageString(registryUrl string, repository string, tag string) string {
	if tag != "" {
		return fmt.Sprintf("%s/%s:%s", strings.Split(registryUrl, "//")[1], repository, tag)
	}
	return fmt.Sprintf("%s/%s", strings.Split(registryUrl, "//")[1], repository)
}

//TODO Communicate with Harbour SCM in order to receive the path to the project-files
func (b Builder) getProjectPath(project string) (string, error) {
	// Just returns a demo path
	return fmt.Sprintf(b.repoPath+"%s", project), nil
}
