package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zot24/repowiki/internal/wiki"
)

var getCmd = &cobra.Command{
	Use:   "get <wiki-path>",
	Short: "Print a wiki document (e.g. `repowiki get pages/overview.md`)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := openWiki()
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(args[0], wiki.DirName+"/")
		p := filepath.Join(w.Root, filepath.FromSlash(rel))
		raw, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("not found: %s (try `repowiki search`)", rel)
		}
		fmt.Print(string(raw))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
