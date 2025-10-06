package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.2.0"

var rootCmd = &cobra.Command{
	Use:   "one",
	Short: "One CLI - Unified tool for freelancers working across multiple projects",
	Long: `One CLI is a unified command-line tool for freelancers working across multiple projects.
It provides a single interface for creating PRs, managing auth, opening tickets, and more.`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("one version %s\n", version))
}
