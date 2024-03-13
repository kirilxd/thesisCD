package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v60/github"
	"io"
	"net/http"
	"os"
	"strings"

	"thesisCD/pkg/kubernetes"
)

type Test struct {
	Test string "json:test"
}

func HandleWebhook(w http.ResponseWriter, r *http.Request, path string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	// Parse the GitHub push event payload
	var pushEvent github.PushEvent
	if err := json.Unmarshal(payload, &pushEvent); err != nil {
		http.Error(w, "Error parsing JSON payload", http.StatusBadRequest)
		return
	}

	// Check for changes in the specified path and apply manifests if necessary
	// This is a simplified example; you'll need to adjust it to your needs
	for _, commit := range pushEvent.Commits {
		for _, modified := range commit.Modified {
			if isPathOfInterest(modified, path) {
				// Fetch file content from GitHub
				content := fetchFileContent(pushEvent.Repo.FullName, modified, commit.ID)
				// Apply the Kubernetes manifest
				fmt.Printf("Applying %s:\n%s\n", modified, content)
				if err := kubernetes.ApplyManifest(content); err != nil {
					fmt.Fprintf(os.Stderr, "Error applying manifest: %v\n", err)
				}
			}
		}
	}

	fmt.Fprintf(w, "Webhook received and processed")
}

func fetchFileContent(repoName *string, filePath string, commitSHA *string) string {
	var ctx = context.Background()
	var githubClient = NewGitHubClient(ctx)

	parts := strings.SplitN(*repoName, "/", 2)
	owner, repo := parts[0], parts[1]
	fileContent, _, _, err := githubClient.Repositories.GetContents(ctx, owner, repo, filePath, &github.RepositoryContentGetOptions{
		Ref: *commitSHA,
	})
	if err != nil {
		fmt.Println("Error fetching file content:", err)
		return ""
	}

	decodedContent, err := base64.StdEncoding.DecodeString(*fileContent.Content)
	if err != nil {
		fmt.Println("Error decoding file content:", err)
		return ""
	}

	return string(decodedContent)
}
