package wiki

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Frontmatter is the OKF-superset metadata block. Unknown keys are preserved
// in Extra so lint --fix and re-serialization never drop user data.
type Frontmatter struct {
	Type        string         `yaml:"type,omitempty"`
	Title       string         `yaml:"title,omitempty"`
	Description string         `yaml:"description,omitempty"`
	Tags        []string       `yaml:"tags,omitempty"`
	Status      string         `yaml:"status,omitempty"`
	Created     string         `yaml:"created,omitempty"`
	Updated     string         `yaml:"updated,omitempty"`
	Confidence  string         `yaml:"confidence,omitempty"`
	OKFVersion  string         `yaml:"okf_version,omitempty"`
	Sources     []SourceRef    `yaml:"sources,omitempty"`
	Generated   *Generated     `yaml:"generated,omitempty"`
	Related     []string       `yaml:"related,omitempty"`
	Extra       map[string]any `yaml:",inline"`
}

// SourceRef is an OKF provenance pointer.
type SourceRef struct {
	Resource string `yaml:"resource"`
	Title    string `yaml:"title,omitempty"`
}

// Generated records which agent/tool produced a document.
type Generated struct {
	By string `yaml:"by,omitempty"`
	At string `yaml:"at,omitempty"`
}

// ValidStatuses are the allowed lifecycle values.
var ValidStatuses = map[string]bool{"draft": true, "stable": true, "deprecated": true}

// Doc is a parsed markdown document: frontmatter + body.
type Doc struct {
	Path string
	FM   *Frontmatter // nil when the file has no frontmatter block
	Body string
}

// ParseFile reads a markdown file and splits frontmatter from body.
func ParseFile(path string) (*Doc, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	doc, err := Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	doc.Path = path
	return doc, nil
}

// Parse splits markdown content into frontmatter and body.
func Parse(raw []byte) (*Doc, error) {
	doc := &Doc{}
	content := string(raw)
	if !strings.HasPrefix(content, "---\n") {
		doc.Body = content
		return doc, nil
	}
	rest := content[len("---\n"):]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		doc.Body = content
		return doc, nil
	}
	yamlPart := rest[:idx+1]
	body := rest[idx+len("\n---"):]
	body = strings.TrimPrefix(body, "\n")
	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(yamlPart), &fm); err != nil {
		return nil, fmt.Errorf("invalid YAML frontmatter: %w", err)
	}
	doc.FM = &fm
	doc.Body = body
	return doc, nil
}

// Render serializes the document back to markdown.
func (d *Doc) Render() ([]byte, error) {
	if d.FM == nil {
		return []byte(d.Body), nil
	}
	var buf bytes.Buffer
	buf.WriteString("---\n")
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(d.FM); err != nil {
		return nil, err
	}
	enc.Close()
	buf.WriteString("---\n\n")
	buf.WriteString(strings.TrimLeft(d.Body, "\n"))
	return buf.Bytes(), nil
}

// Save writes the document back to its path.
func (d *Doc) Save() error {
	out, err := d.Render()
	if err != nil {
		return err
	}
	return os.WriteFile(d.Path, out, 0o644)
}
