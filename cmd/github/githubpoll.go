package github

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"thesisCD/pkg/github"
	"time"
)

var (
	minutesPollInterval int
)

func NewCmdGithubPoll() *cobra.Command {
	var githubPollCmd = &cobra.Command{
		Use:   "githubpoll",
		Short: "Poll updates from github repository on an interval",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting to monitor GitHub for new pushes...")

			ticker := time.NewTicker(time.Duration(minutesPollInterval) * time.Minute)
			quit := make(chan struct{})
			go func() {
				for {
					select {
					case <-ticker.C:
						fmt.Println("Checking for new pushes...")
						githubClient := github.NewGitHubClient(context.Background())
						github.GetCommits(context.Background(), "kirilxd", "thesisCD-infra", "test", githubClient, minutesPollInterval)
					case <-quit:
						ticker.Stop()
						return
					}
				}
			}()

			// Listen for interrupt signal to gracefully shut down
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c
			close(quit)
			fmt.Println("Stopped monitoring GitHub.")
		},
	}

	githubPollCmd.Flags().IntVarP(&minutesPollInterval, "interval", "i", 5, "Poll interval in minutes")

	return githubPollCmd
}
