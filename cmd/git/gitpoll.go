package git

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"thesisCD/pkg/git"
	"time"
)

var (
	secondsPollInterval int
)

func NewCmdGitPoll() *cobra.Command {
	var gitPollCmd = &cobra.Command{
		Use:   "gitpoll",
		Short: "Poll updates from git repository on an interval",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting to monitor repository for updates...")

			path, _ := cmd.Flags().GetString("path")
			repoUrl, _ := cmd.Flags().GetString("repoUrl")

			ticker := time.NewTicker(time.Duration(secondsPollInterval) * time.Second)
			quit := make(chan struct{})
			repo, err := git.CloneRepo(repoUrl)
			if err != nil {
				fmt.Printf("Failed to get repo: %v\n", err)
				return
			}

			go func() {
				for {
					select {
					case <-ticker.C:
						fmt.Println("Checking for updates...")
						var err error
						err = git.PullAndApplyChanges(repo, path)
						if err != nil {
							fmt.Printf("Error pulling updates: %v\n", err)
							// Handle error, could break or log
						}
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
			fmt.Println("Stopped monitoring repository.")
		},
	}

	gitPollCmd.Flags().IntVarP(&secondsPollInterval, "interval", "i", 10, "Poll interval in seconds")

	return gitPollCmd
}
