package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Aliases: []string{"query"},
	Short:   "Full-text search across the wiki",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := openWiki()
		if err != nil {
			return err
		}
		query := strings.ToLower(strings.Join(args, " "))
		files, err := w.MarkdownFiles()
		if err != nil {
			return err
		}
		hits := 0
		for _, p := range files {
			raw, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			for i, line := range strings.Split(string(raw), "\n") {
				if strings.Contains(strings.ToLower(line), query) {
					fmt.Printf("%s:%d: %s\n", w.Rel(p), i+1, strings.TrimSpace(line))
					hits++
				}
			}
		}
		if hits == 0 {
			fmt.Println("No matches.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
