// Command sitebuild assembles the static Carrel website from site/src into site/dist.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// cssOrder is the fixed concatenation order for dist/site.css.
var cssOrder = []string{
	"fonts.css", "tokens.css", "base.css", "layout.css", "chrome.css",
	"hero.css", "window.css",
	"why.css", "problem.css",
	"how.css", "submit.css",
	"sheets.css", "yours.css",
	"themes.css", "faq.css", "cta.css",
}

var defaults = map[string]string{
	"domain":    "[your-domain]",
	"licence":   "[licence]",
	"repo":      "https://github.com/[owner]/carrel",
	"repo_slug": "[owner]/carrel",
	"site_url":  "https://[your-domain]",
}

var includeRe = regexp.MustCompile(`^\s*<!--\s*include:\s*(\S+)\s*-->\s*$`)

const maxDepth = 3

type options struct {
	src, out, scripts string
	release           bool
}

func main() {
	var o options
	serve := flag.String("serve", "", "serve dist on this address and rebuild on change, e.g. 127.0.0.1:8080")
	flag.StringVar(&o.src, "src", "site/src", "source folder")
	flag.StringVar(&o.out, "out", "site/dist", "output folder")
	flag.StringVar(&o.scripts, "scripts", "scripts", "folder holding install.sh and install.ps1")
	flag.BoolVar(&o.release, "release", false, "fail if any placeholder value is left")
	flag.Parse()

	if err := build(o); err != nil {
		log.Fatal(err)
	}
	fmt.Println("built", o.out)
	if *serve == "" {
		return
	}
	go watch(o)
	addr := *serve
	if strings.HasPrefix(addr, ":") {
		addr = "127.0.0.1" + addr
	}
	fmt.Println("serving on http://" + addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(o.out))))
}

func build(o options) error {
	vals, err := loadConfig(filepath.Join(o.src, "site.config.json"))
	if err != nil {
		return err
	}
	if o.release {
		var left []string
		for k, v := range vals {
			if strings.Contains(v, "[") {
				left = append(left, k)
			}
		}
		if len(left) > 0 {
			sort.Strings(left)
			return fmt.Errorf("release: placeholders left in site.config.json: %s", strings.Join(left, ", "))
		}
	}
	if err := os.RemoveAll(o.out); err != nil {
		return err
	}
	if err := os.MkdirAll(o.out, 0o755); err != nil {
		return err
	}
	for _, page := range []string{"index.html", "404.html"} {
		p := filepath.Join(o.src, page)
		if _, err := os.Stat(p); errors.Is(err, fs.ErrNotExist) && page == "404.html" {
			continue
		}
		html, err := resolve(o.src, page, 0)
		if err != nil {
			return err
		}
		if err := writeFile(filepath.Join(o.out, page), replaceTokens(html, vals)); err != nil {
			return err
		}
	}
	var css strings.Builder
	for _, name := range cssOrder {
		b, err := os.ReadFile(filepath.Join(o.src, "css", name))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		css.WriteString("/* " + name + " */\n")
		css.Write(b)
		css.WriteString("\n")
	}
	if err := writeFile(filepath.Join(o.out, "site.css"), replaceTokens(css.String(), vals)); err != nil {
		return err
	}
	if err := copyToken(filepath.Join(o.src, "js", "site.js"), filepath.Join(o.out, "site.js"), vals, false); err != nil {
		return err
	}
	for _, s := range []string{"install.sh", "install.ps1"} {
		if err := copyToken(filepath.Join(o.scripts, s), filepath.Join(o.out, s), vals, true); err != nil {
			return err
		}
	}
	return copyPublic(filepath.Join(filepath.Dir(o.src), "public"), o.out, vals)
}

func loadConfig(path string) (map[string]string, error) {
	vals := map[string]string{}
	for k, v := range defaults {
		vals[k] = v
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return vals, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg map[string]string
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	for k, v := range cfg {
		if v != "" {
			vals[k] = v
		}
	}
	return vals, nil
}

func resolve(root, rel string, depth int) (string, error) {
	if depth > maxDepth {
		return "", fmt.Errorf("include %s: nested more than %d deep", rel, maxDepth)
	}
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return "", fmt.Errorf("include %s: %w", rel, err)
	}
	var out strings.Builder
	sc := bufio.NewScanner(strings.NewReader(string(b)))
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if m := includeRe.FindStringSubmatch(line); m != nil {
			inner, err := resolve(root, m[1], depth+1)
			if err != nil {
				return "", err
			}
			out.WriteString(strings.TrimRight(inner, "\n"))
			out.WriteString("\n")
			continue
		}
		out.WriteString(line)
		out.WriteString("\n")
	}
	return out.String(), sc.Err()
}

func replaceTokens(s string, vals map[string]string) string {
	for k, v := range vals {
		s = strings.ReplaceAll(s, "{{"+k+"}}", v)
	}
	return s
}

func tokenFile(name string) bool {
	switch filepath.Ext(name) {
	case ".html", ".css", ".js", ".sh", ".ps1", ".txt", ".xml":
		return true
	}
	return false
}

func copyToken(src, dst string, vals map[string]string, optional bool) error {
	b, err := os.ReadFile(src)
	if optional && errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return writeFile(dst, replaceTokens(string(b), vals))
}

func copyPublic(pub, out string, vals map[string]string) error {
	if _, err := os.Stat(pub); errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return filepath.WalkDir(pub, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(pub, p)
		dst := filepath.Join(out, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if tokenFile(p) {
			b = []byte(replaceTokens(string(b), vals))
		}
		return os.WriteFile(dst, b, 0o644)
	})
}

func writeFile(path, s string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(s), 0o644)
}

// watch polls source file times and rebuilds when anything changes.
func watch(o options) {
	last := stamp(o)
	for range time.Tick(500 * time.Millisecond) {
		now := stamp(o)
		if now == last {
			continue
		}
		last = now
		if err := build(o); err != nil {
			log.Println("build:", err)
			continue
		}
		log.Println("rebuilt")
	}
}

func stamp(o options) string {
	var b strings.Builder
	for _, dir := range []string{o.src, filepath.Join(filepath.Dir(o.src), "public"), o.scripts} {
		filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if info, err := d.Info(); err == nil {
				fmt.Fprintf(&b, "%s %d %d\n", p, info.ModTime().UnixNano(), info.Size())
			}
			return nil
		})
	}
	return b.String()
}
