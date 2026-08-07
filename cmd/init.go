package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/zot24/repowiki/internal/wiki"
)

var initNoInject bool

var initCmd = &cobra.Command{
	Use:   "init [dir]",
	Short: "Create the .llm-wiki bundle and inject agent instructions",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}
		repoRoot, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		w, created, err := wiki.InitBundle(repoRoot)
		if err != nil {
			return err
		}
		if len(created) == 0 {
			fmt.Println("Wiki already initialized — nothing to create.")
		} else {
			fmt.Printf("Initialized %s\n", filepath.Join(repoRoot, wiki.DirName))
			for _, c := range created {
				fmt.Printf("  + %s\n", c)
			}
		}

		if !initNoInject {
			// AGENTS.md is created if missing; CLAUDE.md only updated if present.
			targets := []struct {
				name   string
				create bool
			}{
				{"AGENTS.md", true},
				{"CLAUDE.md", false},
				{".cursorrules", false},
			}
			for _, t := range targets {
				p := filepath.Join(repoRoot, t.name)
				result, err := wiki.InjectAgentsBlock(p, t.create)
				if err != nil {
					return fmt.Errorf("injecting into %s: %w", t.name, err)
				}
				if result != "skipped" {
					fmt.Printf("  %s: %s\n", t.name, result)
				}
			}
		}

		if len(created) > 0 {
			_ = w.AppendLog("init", "wiki bundle created")
		}
		if _, err := os.Stat(filepath.Join(repoRoot, ".git")); err != nil {
			fmt.Println("\nNote: not a git repository — the wiki works but won't be versioned.")
		}
		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&initNoInject, "no-inject", false, "skip AGENTS.md / CLAUDE.md injection")
	rootCmd.AddCommand(initCmd)
}
