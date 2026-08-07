package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zot24/repowiki/internal/wiki"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Regenerate index.md from pages/ and raw/",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := openWiki()
		if err != nil {
			return err
		}

		var b strings.Builder
		b.WriteString("---\n")
		b.WriteString("type: overview\n")
		b.WriteString("title: \"Wiki index\"\n")
		b.WriteString("description: \"Root index of this OKF knowledge bundle\"\n")
		fmt.Fprintf(&b, "okf_version: %q\n", wiki.OKFVersion)
		b.WriteString("status: stable\n")
		b.WriteString("generated:\n  by: \"tool:repowiki\"\n")
		fmt.Fprintf(&b, "  at: %q\n", wiki.Today())
		b.WriteString("---\n\n# Wiki index\n")

		// Pages, grouped by subdirectory.
		pages, err := w.MarkdownFiles("pages")
		if err != nil {
			return err
		}
		groups := map[string][]*wiki.Doc{}
		for _, p := range pages {
			doc, err := wiki.ParseFile(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: %v (skipped)\n", err)
				continue
			}
			rel := w.Rel(p) // e.g. pages/architecture/api.md
			group := "pages"
			if parts := strings.Split(rel, "/"); len(parts) > 2 {
				group = parts[1]
			}
			groups[group] = append(groups[group], doc)
		}
		var groupNames []string
		for g := range groups {
			groupNames = append(groupNames, g)
		}
		sort.Strings(groupNames)

		b.WriteString("\n## Pages\n")
		empty := true
		for _, g := range groupNames {
			label := strings.ToUpper(g[:1]) + g[1:]
			if g == "pages" {
				label = "General"
			}
			fmt.Fprintf(&b, "\n### %s\n\n", label)
			for _, doc := range groups[g] {
				empty = false
				title := filepath.Base(doc.Path)
				desc := ""
				if doc.FM != nil {
					if doc.FM.Title != "" {
						title = doc.FM.Title
					}
					desc = doc.FM.Description
				}
				line := fmt.Sprintf("- [%s](%s)", title, w.Rel(doc.Path))
				if desc != "" {
					line += " — " + desc
				}
				b.WriteString(line + "\n")
			}
		}
		if empty {
			b.WriteString("\n_No pages yet._\n")
		}

		// Raw layer: counts + recent items per subdir.
		b.WriteString("\n## Raw sources\n\n")
		for _, sub := range wiki.RawSubdirs {
			files, err := w.MarkdownFiles(filepath.Join("raw", sub))
			if err != nil {
				return err
			}
			fmt.Fprintf(&b, "- `raw/%s/` — %d document(s)\n", sub, len(files))
		}

		out := filepath.Join(w.Root, "index.md")
		if err := os.WriteFile(out, []byte(b.String()), 0o644); err != nil {
			return err
		}
		fmt.Printf("Wrote %s (%d pages indexed)\n", w.Rel(out), len(pages))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
