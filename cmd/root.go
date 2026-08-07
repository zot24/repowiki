// Package cmd implements the repowiki CLI commands.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/zot24/repowiki/internal/wiki"
)

var rootCmd = &cobra.Command{
	Use:           "repowiki",
	Short:         "Persistent, Git-native project memory for any coding agent",
	Long:          "repowiki manages a per-repository OKF-compatible LLM wiki (.llm-wiki/)\nthat lives with the code and compounds knowledge across sessions and tools.",
	SilenceUsage:  true,
	SilenceErrors: false,
	Version:       "0.1.0",
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}

// openWiki finds the wiki from the current directory.
func openWiki() (*wiki.Wiki, error) {
	return wiki.Find(".")
}
