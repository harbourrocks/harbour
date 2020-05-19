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
	l "github.com/sirupsen/logrus"
	"io"
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
	redisClient := redisconfig.OpenClient(b.redisOptions)
	redisClient.HSet(job.BuildKey, "build_status", "Running")

	buildCtx, err := b.createBuildContext(job.Request.Project)
	if err != nil {
		l.WithError(err).Error("Failed to create build context")
		return
	}

	opt := types.ImageBuildOptions{
		Tags:       job.Request.Tags,
		Dockerfile: job.Request.Dockerfile,
	}

	resp, err := b.cli.ImageBuild(b.ctx, buildCtx, opt)
	if err != nil {
		l.WithError(err).Error("Failed to build image")
		return
	}

	defer b.cleanUpAfterBuild(buildCtx, resp.Body)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		l.WithError(err).Error("Failed to parse response body")
	}

	logs := buf.String()
	if strings.Contains(logs, "errorDetail") {
		redisClient.HSet(job.BuildKey, "build_status", "Failed", "logs", logs)
		return
	}

	redisClient.HSet(job.BuildKey, "build_status", "Success", "logs", logs)
	l.Tracef("Image %s was built", job.Request.Project)

	imageString := b.getImageString(job.RegistryUrl, job.Request.Project)
	if len(job.Request.Tags) > 0 {
		imageString += ":" + job.Request.Tags[0]
	}

	if err = b.pushImage(imageString, job.RegistryToken); err != nil {
		l.WithError(err).Error("Error while pushing image to registry")
		return
	}

	l.Tracef("Image %s was pushed to registry %s", job.Request.Project, job.RegistryUrl)
}

func (b Builder) cleanUpAfterBuild(buildContext *os.File, logs io.ReadCloser) {
	err := buildContext.Close()
	err = os.Remove(buildContext.Name())
	err = logs.Close()
	if err != nil {
		l.WithError(err).Error("Error while cleaning up build context")
		return
	}
}

func (b Builder) createBuildContext(project string) (*os.File, error) {
	buildContext := fmt.Sprintf(b.ctxPath+"%s.tar", project)
	projectPath, err := b.getProjectPath(project)
	if err != nil {
		l.WithError(err).Error("Failed to receive the project files")
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
			l.WithError(err).Error("Error while closing file")
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

func (b Builder) getImageString(registryUrl string, image string) string {
	return fmt.Sprintf("%s/%s", strings.Split(registryUrl, "//")[1], image)
}

//TODO Communicate with Harbour SCM in order to receive the path to the project-files
func (b Builder) getProjectPath(project string) (string, error) {
	// Just returns a demo path
	return fmt.Sprintf(b.repoPath+"%s", project), nil
}
