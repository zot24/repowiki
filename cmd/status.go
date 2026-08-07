package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zot24/repowiki/internal/wiki"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show wiki overview: counts, inbox, recent activity",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := openWiki()
		if err != nil {
			return err
		}
		fmt.Printf("Wiki: %s\n\n", w.Root)

		pages, err := w.MarkdownFiles("pages")
		if err != nil {
			return err
		}
		fmt.Printf("Pages: %d\n", len(pages))
		for _, sub := range wiki.RawSubdirs {
			files, err := w.MarkdownFiles(filepath.Join("raw", sub))
			if err != nil {
				return err
			}
			fmt.Printf("Raw %-10s %d\n", sub+":", len(files))
		}

		inbox, err := os.ReadDir(filepath.Join(w.Root, "inbox"))
		if err == nil {
			var pending []string
			for _, e := range inbox {
				if !e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
					pending = append(pending, e.Name())
				}
			}
			fmt.Printf("Inbox: %d pending\n", len(pending))
			for _, name := range pending {
				fmt.Printf("  ! inbox/%s (ingest with `repowiki ingest %s`)\n", name, filepath.Join(wiki.DirName, "inbox", name))
			}
		}

		fmt.Println("\nRecent activity:")
		entries := readLogEntries(w)
		if len(entries) == 0 {
			fmt.Println("  (none)")
			return nil
		}
		start := len(entries) - 5
		if start < 0 {
			start = 0
		}
		for _, e := range entries[start:] {
			fmt.Println("  " + e)
		}
		return nil
	},
}

// readLogEntries returns the bullet lines of log.md, oldest first.
func readLogEntries(w *wiki.Wiki) []string {
	raw, err := os.ReadFile(filepath.Join(w.Root, "log.md"))
	if err != nil {
		return nil
	}
	var entries []string
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(line, "- ") {
			entries = append(entries, line)
		}
	}
	return entries
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
