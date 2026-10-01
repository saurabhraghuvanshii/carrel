package server

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

// The UI must work with no network: nothing loaded from another site, and
// every file the CSS or HTML points to is embedded.
func TestWebAssetsAreOffline(t *testing.T) {
	read := func(name string) string {
		b, err := fs.ReadFile(webFS, "web/"+name)
		if err != nil {
			t.Fatalf("missing web/%s: %v", name, err)
		}
		return string(b)
	}
	for _, name := range []string{"index.html", "style.css", "app.js"} {
		if s := read(name); strings.Contains(s, "http://") || strings.Contains(s, "https://") {
			t.Errorf("web/%s refers to the network", name)
		}
	}

	read("vendor/editor.js")
	if !strings.Contains(read("index.html"), `src="vendor/editor.js"`) {
		t.Error("index.html does not load vendor/editor.js")
	}

	for _, icon := range []string{"favicon.ico", "icon.svg"} {
		read(icon)
		if !strings.Contains(read("index.html"), `href="`+icon+`"`) {
			t.Errorf("index.html does not link %s", icon)
		}
	}

	urls := regexp.MustCompile(`url\("([^"]+)"\)`).FindAllStringSubmatch(read("style.css"), -1)
	if len(urls) == 0 {
		t.Fatal("style.css names no font files")
	}
	for _, m := range urls {
		read(m[1])
	}
}

func TestWebAssetsStaySmall(t *testing.T) {
	const limit = 1500 << 10
	total := 0
	err := fs.WalkDir(webFS, "web", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		total += int(info.Size())
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if total > limit {
		t.Fatalf("web/ is %d KB, over the %d KB limit", total>>10, limit>>10)
	}
}
