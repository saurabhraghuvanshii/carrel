package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path, s string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) options {
	t.Helper()
	dir := t.TempDir()
	src := filepath.Join(dir, "site", "src")
	write(t, filepath.Join(src, "index.html"), "<html>\n<!-- include: partials/a.html -->\n</html>\n")
	write(t, filepath.Join(src, "partials", "a.html"), "<p>{{domain}}</p>\n<!-- include: partials/b.html -->\n")
	write(t, filepath.Join(src, "partials", "b.html"), "<a href=\"{{repo}}\">x</a>\n")
	write(t, filepath.Join(src, "css", "base.css"), "b{}\n")
	write(t, filepath.Join(src, "css", "tokens.css"), "t{}\n")
	write(t, filepath.Join(src, "js", "site.js"), "x()\n")
	write(t, filepath.Join(dir, "site", "public", "robots.txt"), "Sitemap: {{site_url}}/sitemap.xml\n")
	write(t, filepath.Join(dir, "scripts", "install.sh"), "REPO=\"{{repo_slug}}\"\n")
	return options{src: src, out: filepath.Join(dir, "site", "dist"), scripts: filepath.Join(dir, "scripts")}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestIncludesAndTokens(t *testing.T) {
	o := fixture(t)
	if err := build(o); err != nil {
		t.Fatal(err)
	}
	got := read(t, filepath.Join(o.out, "index.html"))
	want := "<html>\n<p>[your-domain]</p>\n<a href=\"https://github.com/[owner]/carrel\">x</a>\n</html>\n"
	if got != want {
		t.Fatalf("index.html:\n%q\nwant\n%q", got, want)
	}
	if s := read(t, filepath.Join(o.out, "install.sh")); !strings.Contains(s, `"[owner]/carrel"`) {
		t.Fatalf("install.sh tokens not replaced: %q", s)
	}
	if s := read(t, filepath.Join(o.out, "robots.txt")); !strings.Contains(s, "https://[your-domain]/sitemap.xml") {
		t.Fatalf("robots.txt tokens not replaced: %q", s)
	}
	css := read(t, filepath.Join(o.out, "site.css"))
	if strings.Index(css, "t{}") > strings.Index(css, "b{}") {
		t.Fatalf("css order wrong: %q", css)
	}
}

func TestConfigValues(t *testing.T) {
	o := fixture(t)
	write(t, filepath.Join(o.src, "site.config.json"), `{"domain":"carrel.dev","repo":"https://github.com/me/carrel"}`)
	if err := build(o); err != nil {
		t.Fatal(err)
	}
	got := read(t, filepath.Join(o.out, "index.html"))
	if !strings.Contains(got, "carrel.dev") || !strings.Contains(got, "github.com/me/carrel") {
		t.Fatalf("config values not used: %q", got)
	}
}

func TestMissingIncludeFails(t *testing.T) {
	o := fixture(t)
	write(t, filepath.Join(o.src, "partials", "b.html"), "<!-- include: partials/nope.html -->\n")
	if err := build(o); err == nil || !strings.Contains(err.Error(), "nope.html") {
		t.Fatalf("want missing include error, got %v", err)
	}
}

func TestIncludeDepthLimit(t *testing.T) {
	o := fixture(t)
	write(t, filepath.Join(o.src, "partials", "b.html"), "<!-- include: partials/b.html -->\n")
	if err := build(o); err == nil || !strings.Contains(err.Error(), "deep") {
		t.Fatalf("want depth error, got %v", err)
	}
}

func TestReleaseFailsOnPlaceholders(t *testing.T) {
	o := fixture(t)
	o.release = true
	if err := build(o); err == nil || !strings.Contains(err.Error(), "domain") {
		t.Fatalf("want placeholder error, got %v", err)
	}
	write(t, filepath.Join(o.src, "site.config.json"), `{"domain":"carrel.dev","licence":"MIT","repo":"https://github.com/me/carrel","repo_slug":"me/carrel","site_url":"https://carrel.dev"}`)
	if err := build(o); err != nil {
		t.Fatalf("release with real values: %v", err)
	}
}

func TestRealSiteHasEverySection(t *testing.T) {
	src := filepath.Join("..", "..", "site", "src")
	if _, err := os.Stat(filepath.Join(src, "index.html")); err != nil {
		t.Skip("site/src not present")
	}
	o := options{src: src, out: filepath.Join(t.TempDir(), "dist"), scripts: filepath.Join("..", "..", "scripts")}
	if err := build(o); err != nil {
		t.Fatal(err)
	}
	html := read(t, filepath.Join(o.out, "index.html"))
	for _, id := range []string{"top", "why", "problem", "how", "submit", "sheets", "yours", "themes", "faq"} {
		if !strings.Contains(html, `id="`+id+`"`) {
			t.Errorf("built page has no id=%q", id)
		}
	}
	if strings.Contains(html, "<!-- include:") {
		t.Error("unresolved include left in page")
	}
	if strings.Contains(html, "{{") {
		t.Error("unreplaced token left in page")
	}
}
