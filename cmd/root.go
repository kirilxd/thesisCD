/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"thesisCD/cmd/github"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "thesisCD",
	Short: "Comprehensive science-based GitOps",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(github.NewCmdGithubPoll())
	rootCmd.AddCommand(github.NewCmdGithubWebhook())

	rootCmd.PersistentFlags().String("repo", "", "Repository name")
	rootCmd.PersistentFlags().String("path", "", "Path inside repository")
}
