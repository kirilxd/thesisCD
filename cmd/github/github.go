package github

import (
	"context"
	"github.com/spf13/cobra"
	"thesisCD/pkg/github"
)

func NewCmdGithub() *cobra.Command {
	var githubCmd = &cobra.Command{
		Use:   "github",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			githubClient := github.NewGitHubClient(context.Background())
			github.GetCommits(context.Background(), "kirilxd", "thesisCD-infra", "test", githubClient)
		},
	}

	return githubCmd
}
