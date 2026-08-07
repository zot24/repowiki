package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/zot24/repowiki/internal/wiki"
)

var (
	ingestType  string
	ingestTitle string
	ingestTags  []string
)

var ingestCmd = &cobra.Command{
	Use:   "ingest <file|url|->",
	Short: "Capture a file, URL, or stdin into raw/ with provenance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		w, err := openWiki()
		if err != nil {
			return err
		}
		src := args[0]

		var content, origin string
		switch {
		case src == "-":
			raw, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			content, origin = string(raw), "stdin"
		case strings.HasPrefix(src, "http://"), strings.HasPrefix(src, "https://"):
			client := &http.Client{Timeout: 30 * time.Second}
			resp, err := client.Get(src)
			if err != nil {
				return fmt.Errorf("fetching %s: %w", src, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("fetching %s: HTTP %d", src, resp.StatusCode)
			}
			raw, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
			if err != nil {
				return err
			}
			content, origin = string(raw), src
			if ingestType == "auto" {
				ingestType = "source"
			}
		default:
			raw, err := os.ReadFile(src)
			if err != nil {
				return err
			}
			content, origin = string(raw), src
		}

		// Resolve type: explicit flag > frontmatter in the content > note.
		doc := &wiki.Doc{Body: content}
		if parsed, perr := wiki.Parse([]byte(content)); perr == nil {
			doc = parsed
		}
		typ := ingestType
		if typ == "auto" || typ == "" {
			if doc.FM != nil && doc.FM.Type != "" {
				typ = doc.FM.Type
			} else {
				typ = "note"
			}
		}
		subdir, ok := wiki.TypeDirs[typ]
		if !ok {
			return fmt.Errorf("unknown type %q (want session|decision|note|source|auto)", typ)
		}

		// Resolve title: flag > frontmatter > first heading > filename.
		title := ingestTitle
		if title == "" && doc.FM != nil {
			title = doc.FM.Title
		}
		if title == "" {
			title = firstHeading(doc.Body)
		}
		if title == "" && origin != "stdin" {
			title = strings.TrimSuffix(filepath.Base(origin), filepath.Ext(origin))
		}
		if title == "" {
			title = "untitled"
		}

		if doc.FM == nil {
			doc.FM = &wiki.Frontmatter{}
		}
		doc.FM.Type = typ
		doc.FM.Title = title
		if len(ingestTags) > 0 {
			doc.FM.Tags = append(doc.FM.Tags, ingestTags...)
		}
		if doc.FM.Status == "" {
			doc.FM.Status = "draft"
		}
		if doc.FM.Created == "" {
			doc.FM.Created = wiki.Today()
		}
		doc.FM.Updated = wiki.Today()
		if origin != "stdin" && doc.FM.Generated == nil {
			doc.FM.Sources = append(doc.FM.Sources, wiki.SourceRef{Resource: origin, Title: "original"})
		}

		slug := wiki.Slugify(title)
		if slug == "" {
			slug = "untitled"
		}
		dest := filepath.Join(w.Root, "raw", subdir, fmt.Sprintf("%s-%s.md", wiki.Today(), slug))
		for i := 2; ; i++ {
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				break
			}
			dest = filepath.Join(w.Root, "raw", subdir, fmt.Sprintf("%s-%s-%d.md", wiki.Today(), slug, i))
		}
		doc.Path = dest
		if err := doc.Save(); err != nil {
			return err
		}
		rel := w.Rel(dest)
		fmt.Printf("Ingested %s → %s (type=%s)\n", origin, rel, typ)
		return w.AppendLog("ingest", fmt.Sprintf("%s — %q", rel, title))
	},
}

func firstHeading(body string) string {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(strings.TrimLeft(line, "# "))
		}
	}
	return ""
}

func init() {
	ingestCmd.Flags().StringVar(&ingestType, "type", "auto", "session|decision|note|source|auto")
	ingestCmd.Flags().StringVar(&ingestTitle, "title", "", "override document title")
	ingestCmd.Flags().StringSliceVar(&ingestTags, "tags", nil, "comma-separated tags")
	rootCmd.AddCommand(ingestCmd)
}
