// Command packcheck proves that problem packs are correct. It checks each
// pack's structure, then runs the Java and C++ reference solutions from
// tools/refs through the real runner on the examples, the edge cases and
// random cases from a fixed seed. Every case must pass in both languages.
//
//	go run ./cmd/packcheck [-packs dir] [-refs dir] [problem-id ...]
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"carrel/internal/problems"
	"carrel/internal/runner"
	"carrel/internal/testgen"
)

var refFiles = []struct{ lang, file string }{
	{"java", "Solution.java"},
	{"cpp", "solution.cpp"},
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout))
}

func run(args []string, out io.Writer) int {
	fl := flag.NewFlagSet("packcheck", flag.ContinueOnError)
	fl.SetOutput(out)
	packsDir := fl.String("packs", "internal/problems/packs", "folder with the problem packs")
	refsDir := fl.String("refs", "tools/refs", "folder with the reference solutions")
	random := fl.Int("random", 200, "random cases per pack")
	seed := fl.Int64("seed", 1, "seed for the random cases")
	jobs := fl.Int("jobs", min(runtime.NumCPU(), 8), "packs checked at the same time")
	if err := fl.Parse(args); err != nil {
		return 2
	}

	missing := false
	for _, t := range runner.Doctor() {
		if !t.Found {
			fmt.Fprintf(out, "%s was not found on your PATH. Install it to check packs.\n", t.Name)
			missing = true
		}
	}
	if missing {
		return 2
	}

	lib, err := problems.Load(os.DirFS(*packsDir))
	if err != nil {
		fmt.Fprintf(out, "cannot load packs: %v\n", err)
		return 1
	}
	if len(lib.List()) == 0 {
		fmt.Fprintf(out, "no packs found in %s\n", *packsDir)
		return 1
	}

	ok := true
	if err := problems.Validate(lib); err != nil {
		fmt.Fprintf(out, "structure problems:\n%v\n", err)
		ok = false
	}

	selected := lib.List()
	if ids := fl.Args(); len(ids) > 0 {
		selected = nil
		for _, id := range ids {
			p, found := lib.Get(id)
			if !found {
				fmt.Fprintf(out, "no pack named %q\n", id)
				ok = false
				continue
			}
			selected = append(selected, p)
		}
	}

	// Packs are checked in parallel; each writes to its own buffer, and the
	// buffers are printed in order so the report reads the same every time.
	reports := make([]bytes.Buffer, len(selected))
	passed := make([]bool, len(selected))
	next := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < max(1, *jobs); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range next {
				passed[i] = checkPack(&reports[i], selected[i], *refsDir, *seed, *random)
			}
		}()
	}
	for i := range selected {
		next <- i
	}
	close(next)
	wg.Wait()
	for i := range selected {
		_, _ = out.Write(reports[i].Bytes())
		ok = ok && passed[i]
	}

	if !ok {
		fmt.Fprintln(out, "pack check failed")
		return 1
	}
	fmt.Fprintf(out, "all %d packs passed\n", len(selected))
	return 0
}

func checkPack(out io.Writer, p *problems.Problem, refsDir string, seed int64, random int) bool {
	cases, err := casesFor(p, seed, random)
	if err != nil {
		fmt.Fprintf(out, "FAIL  %s: %v\n", p.ID, err)
		return false
	}
	ok := true
	for _, rf := range refFiles {
		if !checkLang(out, p, rf.lang, filepath.Join(refsDir, p.ID, rf.file), cases) {
			ok = false
		}
	}
	return ok
}

func casesFor(p *problems.Problem, seed int64, random int) ([]runner.Case, error) {
	var cases []runner.Case
	for _, c := range p.Examples {
		cases = append(cases, runner.Case{Kind: "example", Label: c.Label, Input: c.Input, Expected: c.Expected})
	}
	for _, c := range p.Edge {
		cases = append(cases, runner.Case{Kind: "edge", Label: c.Label, Input: c.Input, Expected: c.Expected})
	}
	gen, err := testgen.Generate(p.Generator, seed, random)
	if err != nil {
		return nil, err
	}
	for i, c := range gen {
		cases = append(cases, runner.Case{Kind: "random", Label: "Random " + strconv.Itoa(i+1), Input: c.Input, Expected: c.Expected})
	}
	return cases, nil
}

func checkLang(out io.Writer, p *problems.Problem, lang, refPath string, cases []runner.Case) bool {
	code, err := os.ReadFile(refPath)
	if err != nil {
		fmt.Fprintf(out, "FAIL  %-28s %-4s no reference at %s\n", p.ID, lang, refPath)
		return false
	}
	driver, err := p.File("driver." + lang)
	if err != nil {
		fmt.Fprintf(out, "FAIL  %-28s %-4s no driver.%s in the pack\n", p.ID, lang, lang)
		return false
	}
	rep := runner.Run(context.Background(), runner.Request{
		Lang: lang, Code: string(code), Driver: driver, Cases: cases,
		CompileTimeout: 60 * time.Second, RunTimeout: 30 * time.Second,
	})
	passed := rep.Status == "ok" && rep.Passed == rep.Total
	word := "ok  "
	if !passed {
		word = "FAIL"
	}
	fmt.Fprintf(out, "%s  %-28s %-4s %d/%d  %dms\n", word, p.ID, lang, rep.Passed, rep.Total, rep.DurationMs)
	if passed {
		return true
	}
	if rep.Status != "ok" {
		fmt.Fprintf(out, "      status %s\n%s\n", rep.Status, indent(rep.Message))
	}
	for _, r := range rep.Results {
		if !r.Passed {
			fmt.Fprintf(out, "      first failing case: %s\n      input:\n%s\n      expected: %s\n      got:      %s\n",
				r.Label, indent(r.Input), r.Expected, r.Got)
			break
		}
	}
	return false
}

func indent(s string) string {
	if s == "" {
		return "        (empty)"
	}
	return "        " + strings.ReplaceAll(s, "\n", "\n        ")
}
