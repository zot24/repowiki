package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zot24/repowiki/internal/wiki"
)

var lintFix bool

var mdLinkRe = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Check wiki structure, frontmatter, and OKF conformance",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := openWiki()
		if err != nil {
			return err
		}
		var problems []string
		report := func(format string, a ...any) {
			problems = append(problems, fmt.Sprintf(format, a...))
		}

		// Root index must exist and declare okf_version.
		indexPath := filepath.Join(w.Root, "index.md")
		if doc, err := wiki.ParseFile(indexPath); err != nil {
			report("index.md: %v", err)
		} else if doc.FM == nil || doc.FM.OKFVersion == "" {
			report("index.md: missing okf_version in frontmatter (OKF requires it at the bundle root)")
		}

		files, err := w.MarkdownFiles("pages", "raw")
		if err != nil {
			return err
		}
		for _, p := range files {
			rel := w.Rel(p)
			doc, err := wiki.ParseFile(p)
			if err != nil {
				report("%s: %v", rel, err)
				continue
			}
			if doc.FM == nil {
				if lintFix {
					doc.FM = &wiki.Frontmatter{
						Type:    guessType(rel),
						Title:   strings.TrimSuffix(filepath.Base(p), ".md"),
						Status:  "draft",
						Created: wiki.Today(),
						Updated: wiki.Today(),
					}
					if err := doc.Save(); err != nil {
						return err
					}
					fmt.Printf("fixed: %s — added frontmatter stub\n", rel)
				} else {
					report("%s: no frontmatter (run `repowiki lint --fix` to add a stub)", rel)
				}
				continue
			}
			if doc.FM.Type == "" {
				report("%s: frontmatter missing required `type`", rel)
			}
			if doc.FM.Title == "" {
				report("%s: frontmatter missing recommended `title`", rel)
			}
			if doc.FM.Status != "" && !wiki.ValidStatuses[doc.FM.Status] {
				report("%s: invalid status %q (want draft|stable|deprecated)", rel, doc.FM.Status)
			}
			// Relative links must resolve.
			for _, m := range mdLinkRe.FindAllStringSubmatch(doc.Body, -1) {
				target := m[1]
				if strings.Contains(target, "://") || strings.HasPrefix(target, "#") || strings.HasPrefix(target, "mailto:") {
					continue
				}
				target = strings.SplitN(target, "#", 2)[0]
				if target == "" {
					continue
				}
				var full string
				if strings.HasPrefix(target, "/") {
					full = filepath.Join(w.Root, target)
				} else {
					full = filepath.Join(filepath.Dir(p), target)
				}
				if _, err := os.Stat(full); err != nil {
					report("%s: broken link → %s", rel, m[1])
				}
			}
		}

		if len(problems) == 0 {
			fmt.Println("OK — no problems found.")
			return nil
		}
		for _, p := range problems {
			fmt.Println("✗ " + p)
		}
		return fmt.Errorf("%d problem(s) found", len(problems))
	},
}

// guessType infers a document type from its wiki-relative path.
func guessType(rel string) string {
	parts := strings.Split(rel, "/")
	if len(parts) >= 2 {
		switch parts[0] {
		case "raw":
			for typ, dir := range wiki.TypeDirs {
				if parts[1] == dir {
					return typ
				}
			}
		case "pages":
			switch parts[1] {
			case "architecture", "decisions", "plans", "components", "concepts":
				return strings.TrimSuffix(parts[1], "s")
			}
		}
	}
	return "note"
}

func init() {
	lintCmd.Flags().BoolVar(&lintFix, "fix", false, "auto-fix what can be fixed")
	rootCmd.AddCommand(lintCmd)
}
