package git

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"os"
)

// CloneRepository clones the specified git repository
//  authentication has to be included in the url
func CloneRepository(ctx context.Context, url, path string) (repository *git.Repository, worktree *git.Worktree, err error) {
	log := logconfig.GetLogCtx(ctx)

	repository, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		log.WithError(err).WithField("repoUrl", url).Error("Failed to clone repository")
		return
	}

	worktree, err = repository.Worktree()
	if err != nil {
		log.WithError(err).WithField("repoUrl", url).Error("Failed to get worktree")
		return
	}

	return
}
