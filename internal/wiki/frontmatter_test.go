package wiki

import (
	"strings"
	"testing"
)

func TestParseRoundTrip(t *testing.T) {
	in := `---
type: decision
title: "Auth choice"
tags: [auth, security]
status: stable
custom_key: kept
---

# Auth choice

Body text.
`
	doc, err := Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if doc.FM.Type != "decision" || doc.FM.Title != "Auth choice" {
		t.Errorf("bad frontmatter: %+v", doc.FM)
	}
	if doc.FM.Extra["custom_key"] != "kept" {
		t.Errorf("unknown key not preserved: %v", doc.FM.Extra)
	}
	out, err := doc.Render()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"type: decision", "custom_key: kept", "# Auth choice", "Body text."} {
		if !strings.Contains(string(out), want) {
			t.Errorf("render missing %q:\n%s", want, out)
		}
	}
}

func TestParseNoFrontmatter(t *testing.T) {
	doc, err := Parse([]byte("just text\n"))
	if err != nil {
		t.Fatal(err)
	}
	if doc.FM != nil {
		t.Errorf("expected nil frontmatter, got %+v", doc.FM)
	}
	if doc.Body != "just text\n" {
		t.Errorf("body mangled: %q", doc.Body)
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Chose SQLite for storage": "chose-sqlite-for-storage",
		"  Weird -- chars!! (v2) ": "weird-chars-v2",
		"":                         "",
	}
	for in, want := range cases {
		if got := Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}
