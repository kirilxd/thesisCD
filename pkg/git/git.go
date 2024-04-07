package git

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
	"strings"
	"thesisCD/pkg/kubernetes"
)

func PullAndApplyChanges(repo *git.Repository, path string) error {
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	oldHead, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD before pull: %v", err)
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull: %v", err)
	} else if err == git.NoErrAlreadyUpToDate {
		fmt.Println("No new changes to pull.")
		return nil
	}

	newHead, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD after pull: %v", err)
	}

	if newHead.Hash() != oldHead.Hash() {
		err := checkPathChanges(repo, oldHead, newHead, path)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("No new commits were merged during the pull.")
	}

	return nil

}

func checkPathChanges(repo *git.Repository, oldHead, newHead *plumbing.Reference, path string) error {
	oldCommit, err := repo.CommitObject(oldHead.Hash())
	if err != nil {
		return fmt.Errorf("failed to get old HEAD commit: %v", err)
	}

	newCommit, err := repo.CommitObject(newHead.Hash())
	if err != nil {
		return fmt.Errorf("failed to get new HEAD commit: %v", err)
	}

	diff, err := oldCommit.Patch(newCommit)
	if err != nil {
		return fmt.Errorf("failed to get diff: %v", err)
	}

	for _, filePatch := range diff.FilePatches() {
		from, to := filePatch.Files()

		if to != nil && strings.Contains(to.Path(), path) {
			content, err := getFileContentAtCommit(newCommit, to.Path())
			if err != nil {
				return fmt.Errorf("failed to get file content for %s: %v", to.Path(), err)
			}
			err = kubernetes.ApplyManifest(content)
			if err != nil {
				return err
			}
		}
		if to == nil && strings.Contains(from.Path(), path) {
			content, err := getFileContentAtCommit(oldCommit, from.Path())
			if err != nil {
				return fmt.Errorf("failed to get file content for %s: %v", from.Path(), err)
			}
			err = kubernetes.DeleteResource(content)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CloneRepo(repoUrl string) (*git.Repository, error) {
	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:           repoUrl,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName("refs/heads/main"),
		SingleBranch:  true,
	})

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to clone repository: %v", err)
	}

	fmt.Println("Repository cloned into memory successfully.")
	return repo, nil
}

func getFileContentAtCommit(commit *object.Commit, filePath string) (string, error) {
	file, err := commit.File(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file %s from commit: %v", filePath, err)
	}

	reader, err := file.Reader()
	if err != nil {
		return "", fmt.Errorf("failed to open file reader for %s: %v", filePath, err)
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read file content for %s: %v", filePath, err)
	}

	return buf.String(), nil
}
