// Package wiki implements the on-disk model of a .llm-wiki OKF bundle.
package wiki

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	// DirName is the wiki bundle directory at the repo root.
	DirName = ".llm-wiki"
	// OKFVersion is the Open Knowledge Format version we target.
	OKFVersion = "0.2"
)

// RawSubdirs are the immutable-source directories under raw/.
var RawSubdirs = []string{"sessions", "notes", "sources", "decisions"}

// PageSubdirs are the compiled-knowledge directories under pages/.
var PageSubdirs = []string{"architecture", "decisions", "plans", "components", "concepts"}

// TypeDirs maps an ingest --type to its raw/ subdirectory.
var TypeDirs = map[string]string{
	"session":  "sessions",
	"note":     "notes",
	"source":   "sources",
	"decision": "decisions",
}

// Wiki is an opened .llm-wiki bundle.
type Wiki struct {
	RepoRoot string // directory containing .llm-wiki
	Root     string // the .llm-wiki directory itself
}

// Find walks upward from dir looking for a .llm-wiki directory.
func Find(dir string) (*Wiki, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for d := abs; ; d = filepath.Dir(d) {
		candidate := filepath.Join(d, DirName)
		if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
			return &Wiki{RepoRoot: d, Root: candidate}, nil
		}
		if d == filepath.Dir(d) {
			return nil, fmt.Errorf("no %s found in %s or any parent — run `repowiki init` first", DirName, abs)
		}
	}
}

// Rel returns p relative to the wiki root, for display and links.
func (w *Wiki) Rel(p string) string {
	rel, err := filepath.Rel(w.Root, p)
	if err != nil {
		return p
	}
	return filepath.ToSlash(rel)
}

// AppendLog appends a timestamped entry to log.md.
func (w *Wiki) AppendLog(action, detail string) error {
	line := fmt.Sprintf("- %s — **%s** — %s\n", time.Now().Format("2006-01-02T15:04-07:00"), action, detail)
	f, err := os.OpenFile(filepath.Join(w.Root, "log.md"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}

// MarkdownFiles returns all .md files under the given wiki subdirectories
// (or the whole wiki if none given), sorted, absolute paths.
func (w *Wiki) MarkdownFiles(subdirs ...string) ([]string, error) {
	roots := []string{w.Root}
	if len(subdirs) > 0 {
		roots = nil
		for _, s := range subdirs {
			roots = append(roots, filepath.Join(w.Root, s))
		}
	}
	var files []string
	for _, root := range roots {
		err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					return filepath.SkipDir
				}
				return err
			}
			if d.IsDir() {
				if strings.HasPrefix(d.Name(), ".") && p != root {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(p, ".md") {
				files = append(files, p)
			}
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	sort.Strings(files)
	return files, nil
}

// Slugify turns a title into a filename-safe slug.
func Slugify(s string) string {
	var b strings.Builder
	prevDash := true
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}
