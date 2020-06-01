package harbourbuild

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/builder/dockerignore"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/fileutils"

	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/jhoonb/archivex"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Builder struct {
	jobChan      chan models.BuildJob
	cli          *client.Client
	ctx          context.Context
	ctxPath      string
	redisOptions redisconfig.RedisOptions
	log          *logrus.Entry
}

func NewBuilder(jobChan chan models.BuildJob, ctxPath string, redisConfig redisconfig.RedisOptions) (Builder, error) {
	var builder Builder
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return builder, err
	}

	builder = Builder{jobChan: jobChan, cli: cli, ctx: ctx, ctxPath: ctxPath, redisOptions: redisConfig}
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
	b.log = logrus.WithField("reqId", job.ReqId)
	redisClient := redisconfig.OpenClient(b.redisOptions)

	if err := redisClient.HSet(job.BuildKey, "build_status", "Running").Err(); err != nil {
		b.log.WithError(err).Error("Failed to save data to redis")
		return
	}

	buildCtx, err := b.createBuildContext(job.FilePath, job.Dockerfile)
	if err != nil {
		b.log.WithError(err).Error("Failed to create build context")
		return
	}

	tag := []string{b.getImageString(job.RegistryUrl, job.Repository, job.Tag)}
	opt := types.ImageBuildOptions{
		Tags:       tag,
		Dockerfile: job.Dockerfile,
	}

	resp, err := b.cli.ImageBuild(b.ctx, buildCtx, opt)
	if err != nil {
		b.log.WithError(err).Error("Failed to build image")
		return
	}

	defer func() {
		err := buildCtx.Close()
		err = os.Remove(buildCtx.Name())
		err = os.RemoveAll(job.FilePath)
		err = resp.Body.Close()
		if err != nil {
			b.log.WithError(err).Error("Error while cleaning up build context")
			return
		}
	}()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		b.log.WithError(err).Error("Failed to parse response body")
	}

	logs := buf.String()
	if strings.Contains(logs, "errorDetail") {
		b.log.Errorf("Build failed: %s", logs)
		if err := redisClient.HSet(job.BuildKey, "build_status", "Failed", "logs", logs).Err(); err != nil {
			b.log.WithError(err).Error("Failed to save data to redis")
			return
		}
		return
	}

	if err := redisClient.HSet(job.BuildKey, "build_status", "Success", "logs", logs).Err(); err != nil {
		b.log.WithError(err).Error("Failed to save data to redis")
		return
	}
	b.log.Tracef("Image %s was built", job.Repository)

	imageString := b.getImageString(job.RegistryUrl, job.Repository, job.Tag)

	if err = b.pushImage(imageString, job.RegistryToken); err != nil {
		b.log.WithError(err).Error("Error while pushing image to registry")
		return
	}

	b.log.Tracef("Image %s was pushed to registry %s", job.Repository, job.RegistryUrl)
}

func (b Builder) createBuildContext(filePath string, dockerfile string) (*os.File, error) {
	var excludes = []string{}
	buildContext := fmt.Sprintf("%s/%s.tar", b.ctxPath, strings.Split(filePath, "/")[1])

	path := fmt.Sprintf("./%s%s%s", filePath, dockerfile, ".dockerignore")
	ignore, err := os.Open(path)
	if os.IsNotExist(err) {
		b.log.Trace("Fallback to root .dockerignore")
		ignore, err = os.Open(fmt.Sprintf("%s/%s", filePath, ".dockerignore"))
		if os.IsNotExist(err) {
			b.log.Trace("No .dockerignore found")
		}
	}

	excludes, err = dockerignore.ReadAll(ignore)
	if err != nil {

	}

	patternMatcher, err := fileutils.NewPatternMatcher(excludes)
	if err != nil {

	}

	tar := new(archivex.TarFile)
	err = tar.Create(buildContext)
	if err != nil {
		logrus.WithError(err).Error("Couldn't create build context")
	}

	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		split := strings.Split(filepath.Clean(path), filePath+"/")

		if len(split) > 1 {
			excluded, _ := patternMatcher.Matches(split[1])
			if !excluded {
				file, _ := os.Open(path)
				tar.Add(split[1], file, nil)
			}
			return nil
		}

		return nil
	})

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
	fmt.Println(logs)
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
