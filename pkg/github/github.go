package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"thesisCD/pkg/kubernetes"
	"time"
)

const MINUTES_COMMIT_CHECK = 60

func NewGitHubClient(ctx context.Context) *github.Client {
	ghOwner := os.Getenv("GITHUB_USERNAME")
	ghToken := os.Getenv("GITHUB_TOKEN")

	if ghOwner == "" {
		fmt.Println("Error: GITHUB_USERNAME environment variable is not set.")
	}

	if ghToken == "" {
		fmt.Println("Error: GITHUB_TOKEN environment variable is not set.")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

func GetCommits(ctx context.Context, owner string, repo string, path string, client *github.Client) {
	commits, _, err := client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
		Since: time.Now().Add(-MINUTES_COMMIT_CHECK * time.Minute),
	})

	if err != nil {
		fmt.Println("Error fetching commits:", err)
		return
	}

	for _, commit := range commits {
		commitDetail, _, err := client.Repositories.GetCommit(ctx, owner, repo, *commit.SHA, nil)
		if err != nil {
			fmt.Println("Error fetching commit details:", err)
			continue
		}

		for _, file := range commitDetail.Files {
			if strings.HasPrefix(file.GetFilename(), path) {
				fmt.Printf("File in the target path changed: %s in commit %s\n", file.GetFilename(), *commit.SHA)

				// Fetch and display the content of the changed file
				content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, file.GetFilename(), &github.RepositoryContentGetOptions{
					Ref: *commit.SHA,
				})
				if err != nil {
					fmt.Println("Error fetching file content:", err)
					continue
				}

				decodedContent, err := base64.StdEncoding.DecodeString(*content.Content)
				if err != nil {
					fmt.Println("Error decoding file content:", err)
					continue
				}

				err = kubernetes.ApplyManifest(string(decodedContent))
				if err != nil {
					continue
				}
				fmt.Printf("Content of %s:\n%s\n", file.GetFilename(), string(decodedContent))
			}
		}

	}

	if len(commits) == 0 {
		fmt.Printf("No new pushes in the last %d minutes.\n", MINUTES_COMMIT_CHECK)
	}
}
