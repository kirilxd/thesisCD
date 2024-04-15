package git

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"thesisCD/pkg/git"
)

func NewCmdGitWebhook() *cobra.Command {
	var gitWebhookCmd = &cobra.Command{
		Use:   "gitpush",
		Short: "Listen for updates in repository with webhook",
		Run: func(cmd *cobra.Command, args []string) {
			path, _ := cmd.Flags().GetString("path")
			repoUrl, _ := cmd.Flags().GetString("repoUrl")
			port, _ := cmd.Flags().GetString("port")

			repo, _ := git.CloneRepo(repoUrl)

			http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
				git.HandleGiteaWebhook(w, r, path, repo)
			})
			fmt.Printf("Starting server to listen for webhook events on port %s\n", port)
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				fmt.Printf("Error starting server: %s\n", err)
			}
		},
	}

	return gitWebhookCmd
}
