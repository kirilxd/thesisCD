package git

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"io"
	"net/http"
	"strings"
)

type GiteaPushEvent struct {
	Commits []struct {
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"commits"`
}

func HandleGiteaWebhook(w http.ResponseWriter, r *http.Request, path string, repo *git.Repository) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	var pushEvent GiteaPushEvent
	if err := json.Unmarshal(payload, &pushEvent); err != nil {
		http.Error(w, "Error parsing JSON payload", http.StatusBadRequest)
		return
	}

	pathAffected := false
	for _, commit := range pushEvent.Commits {
		for _, filePath := range append(commit.Added, append(commit.Removed, commit.Modified...)...) {
			if strings.HasPrefix(filePath, path) {
				pathAffected = true
				break
			}
		}
		if pathAffected {
			break
		}
	}

	if !pathAffected {
		fmt.Println("No changes detected in the specified path. Stopping further processing.")
		return
	}

	err = PullAndApplyChanges(repo, path)
	if err != nil {
		http.Error(w, "Failed to pull and apply changes", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Webhook received and processed")
}
