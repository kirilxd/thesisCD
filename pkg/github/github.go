package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/v60/github"
	"thesisCD/pkg/kubernetes"
	"time"
)

func GetCommits(ctx context.Context, owner string, repo string, path string, client *github.Client, interval int) {
	commits, _, err := client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
		Since: time.Now().Add(-time.Duration(interval) * time.Minute),
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
			if isPathOfInterest(file.GetFilename(), path) {
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
		fmt.Printf("No new pushes in the last %d minutes.\n", interval)
	}
}
