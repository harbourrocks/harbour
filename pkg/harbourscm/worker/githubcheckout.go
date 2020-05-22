package worker

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	hGit "github.com/harbourrocks/harbour/pkg/harbourscm/git"
	"github.com/harbourrocks/harbour/pkg/harbourscm/github"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"os"
	"path"
)

type GithubCheckoutTask struct {
	OrganizationLogin string
	Repository        string
	CallbackUrl       string
	State             string
	Commit            plumbing.Hash
	Ctx               context.Context
}

func CheckoutGithub(task GithubCheckoutTask) {
	ctx := task.Ctx
	scmConfig := configuration.GetSCMConfigCtx(ctx)
	log := logconfig.GetLogCtx(ctx)

	var err error
	var workspacePath string

	log.Tracef("Checkout github repository: %v", task)

	// call callback on success and error
	defer func() {
		callbackModel := CheckoutCompletedModel{}
		if err != nil {
			callbackModel.Success = false
		} else {
			callbackModel.Success = true
			callbackModel.WorkspacePath = workspacePath
		}

		resp, err := apiclient.Post(ctx, fmt.Sprintf("%s?state=%s", task.CallbackUrl, task.State), nil, callbackModel, "", nil)
		if err != nil || resp.StatusCode >= 300 {
			log.WithError(err).Error("Callback failed")
			return
		}
	}()

	// generate a github token
	installationToken, err := github.GenerateTokenForOrganization(ctx, task.OrganizationLogin)
	if err != nil {
		return // error logged in GenerateTokenForOrganization
	}

	// create workspace
	checkoutFolder := uuid.New().String()
	workspacePath = path.Join(scmConfig.CheckoutPath, checkoutFolder)
	err = os.Mkdir(workspacePath, os.ModeDir)
	if err != nil {
		log.WithError(err).WithField("workspacePath", workspacePath).Error("Failed to generate workspace directors")
		return
	} else {
		log.WithField("workspacePath", workspacePath).Trace("Workspace created")
	}

	// build github repository url and checkout repository
	gitUrl := fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git", installationToken, task.OrganizationLogin, task.Repository)
	_, worktree, err := hGit.CloneRepository(ctx, gitUrl, workspacePath)
	if err != nil {
		return // error logged in CloneRepository
	}

	// checkout specific commit
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: task.Commit,
	})
	if err != nil {
		log.WithError(err).WithField("commitHash", task.Commit.String()).Error("Failed to checkout commit")
		return
	}
}
