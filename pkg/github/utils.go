package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

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

func isPathOfInterest(fileName string, path string) bool {
	return strings.HasPrefix(fileName, path)
}
