package harbourbuild

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/jhoonb/archivex"
	l "github.com/sirupsen/logrus"
	"os"
)

type Builder struct {
	jobChan  chan models.BuildJob
	cli      *client.Client
	ctx      context.Context
	ctxPath  string
	repoPath string
}

func NewBuilder(jobChan chan models.BuildJob, ctxPath string, repoPath string) (Builder, error) {
	var builder Builder
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return builder, err
	}

	builder = Builder{jobChan: jobChan, cli: cli, ctx: ctx, ctxPath: ctxPath, repoPath: repoPath}
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
	l.Trace("Create Build Context")
	buildCtx, err := b.createBuildContext(job.Request.Project)
	if err != nil {
		l.WithError(err).Error("Failed to create build context")
		return
	}

	defer b.cleanUpAfterBuild(buildCtx)

	l.Trace("Create Build Options")
	opt := types.ImageBuildOptions{
		Tags:       job.Request.Tags,
		Dockerfile: job.Request.Dockerfile,
	}

	l.Trace("Build Image")
	_, err = b.cli.ImageBuild(b.ctx, buildCtx, opt)
	if err != nil {
		l.WithError(err).Error("Failed to build image")
		return
	}
}

func (b Builder) cleanUpAfterBuild(buildContext *os.File) {
	l.Trace("Cleaning up")
	err := buildContext.Close()
	err = os.Remove(buildContext.Name())
	if err != nil {
		l.WithError(err).Error("Error while cleaning up build context")
		return
	}
}

//TODO Communicate with Harbour SCM in order to receive the path to the project-files
func (b Builder) getProjectPath(project string) (string, error) {
	// Just returns a demo path
	return fmt.Sprintf(b.repoPath+"%s", project), nil
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
