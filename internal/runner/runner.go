// Package runner compiles and runs the user's code against a list of test
// cases using the compilers already installed on this computer.
//
// Test protocol: the driver reads "T" and then T cases from stdin and prints
// exactly one line of output per case. A case that throws prints "ERROR ...".
package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Case struct {
	Kind     string // example | edge | random
	Label    string
	Input    string
	Expected string
}

type Result struct {
	Kind     string `json:"kind"`
	Label    string `json:"label"`
	Passed   bool   `json:"passed"`
	Input    string `json:"input,omitempty"`
	Expected string `json:"expected,omitempty"`
	Got      string `json:"got,omitempty"`
}

type Report struct {
	// Status is ok, compile_error, runtime_error, timeout, tooling_missing or internal_error.
	Status     string   `json:"status"`
	Message    string   `json:"message,omitempty"`
	Results    []Result `json:"results"`
	Passed     int      `json:"passed"`
	Total      int      `json:"total"`
	DurationMs int64    `json:"durationMs"`
}

type Request struct {
	Lang           string
	Code           string
	Driver         string
	Cases          []Case
	CompileTimeout time.Duration
	RunTimeout     time.Duration
}

type spec struct {
	tools      []string
	codeFile   string
	driverFile string
	compile    []string
	run        []string
}

func specFor(lang string) (spec, bool) {
	switch lang {
	case "java":
		return spec{
			tools:      []string{"javac", "java"},
			codeFile:   "Solution.java",
			driverFile: "Main.java",
			compile:    []string{"javac", "-d", "out", "Solution.java", "Main.java"},
			run:        []string{"java", "-Xss64m", "-XX:+UseSerialGC", "-cp", "out", "Main"},
		}, true
	case "cpp":
		exe := "prog"
		if runtime.GOOS == "windows" {
			exe += ".exe"
		}
		return spec{
			tools:      []string{"g++"},
			codeFile:   "solution.cpp",
			driverFile: "driver.cpp",
			compile:    []string{"g++", "-O2", "-std=c++17", "-o", exe, "driver.cpp"},
			run:        []string{"." + string(filepath.Separator) + exe},
		}, true
	}
	return spec{}, false
}

// Run compiles the code once and runs every case in one process.
func Run(ctx context.Context, req Request) Report {
	start := time.Now()
	rep := Report{Results: []Result{}, Total: len(req.Cases)}
	finish := func(status, msg string) Report {
		rep.Status, rep.Message = status, msg
		rep.DurationMs = time.Since(start).Milliseconds()
		return rep
	}

	sp, ok := specFor(req.Lang)
	if !ok {
		return finish("internal_error", "unsupported language: "+req.Lang)
	}
	for _, tool := range sp.tools {
		if _, err := exec.LookPath(tool); err != nil {
			return finish("tooling_missing", tool+" was not found on your PATH. Install it, then run `dsa doctor`.")
		}
	}

	dir, err := os.MkdirTemp("", "dsa-run-*")
	if err != nil {
		return finish("internal_error", err.Error())
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(filepath.Join(dir, sp.codeFile), []byte(req.Code), 0o600); err != nil {
		return finish("internal_error", err.Error())
	}
	if err := os.WriteFile(filepath.Join(dir, sp.driverFile), []byte(req.Driver), 0o600); err != nil {
		return finish("internal_error", err.Error())
	}

	cctx, cancel := context.WithTimeout(ctx, req.CompileTimeout)
	defer cancel()
	stdout, stderr, err := execute(cctx, dir, nil, sp.compile[0], sp.compile[1:]...)
	if err != nil {
		if errors.Is(cctx.Err(), context.DeadlineExceeded) {
			return finish("timeout", "Compiling took too long.")
		}
		return finish("compile_error", clean(stdout+stderr))
	}

	var stdin bytes.Buffer
	fmt.Fprintf(&stdin, "%d\n", len(req.Cases))
	for _, c := range req.Cases {
		stdin.WriteString(c.Input)
		stdin.WriteString("\n")
	}

	rctx, rcancel := context.WithTimeout(ctx, req.RunTimeout)
	defer rcancel()
	stdout, stderr, runErr := execute(rctx, dir, &stdin, sp.run[0], sp.run[1:]...)

	var lines []string
	if s := strings.TrimRight(stdout, "\r\n"); s != "" {
		lines = strings.Split(s, "\n")
	}
	for i, c := range req.Cases {
		got := ""
		have := i < len(lines)
		if have {
			got = strings.TrimSpace(lines[i])
		}
		passed := have && got == strings.TrimSpace(c.Expected)
		r := Result{Kind: c.Kind, Label: c.Label, Passed: passed}
		if passed {
			rep.Passed++
		}
		if c.Kind == "example" || !passed {
			r.Input, r.Expected, r.Got = shorten(c.Input), shorten(c.Expected), shorten(got)
		}
		rep.Results = append(rep.Results, r)
	}

	switch {
	case errors.Is(rctx.Err(), context.DeadlineExceeded):
		return finish("timeout", fmt.Sprintf("Your solution took longer than %s.", req.RunTimeout))
	case runErr != nil:
		return finish("runtime_error", clean(stderr))
	case len(lines) < len(req.Cases):
		return finish("runtime_error", "The program stopped before every case ran.\n"+clean(stderr))
	}
	return finish("ok", "")
}

func execute(ctx context.Context, dir string, stdin *bytes.Buffer, name string, args ...string) (string, string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.WaitDelay = 2 * time.Second
	if stdin != nil {
		cmd.Stdin = stdin
	}
	out, errb := &limitedBuffer{max: 1 << 20}, &limitedBuffer{max: 64 << 10}
	cmd.Stdout, cmd.Stderr = out, errb
	err := cmd.Run()
	return out.buf.String(), errb.buf.String(), err
}

type limitedBuffer struct {
	buf bytes.Buffer
	max int
}

func (l *limitedBuffer) Write(p []byte) (int, error) {
	if room := l.max - l.buf.Len(); room > 0 {
		if len(p) > room {
			l.buf.Write(p[:room])
		} else {
			l.buf.Write(p)
		}
	}
	return len(p), nil
}

// clean drops launcher noise and keeps messages a reasonable length.
func clean(s string) string {
	var keep []string
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "Picked up ") {
			continue
		}
		keep = append(keep, line)
	}
	return shorten2(strings.TrimSpace(strings.Join(keep, "\n")), 2000)
}

func shorten(s string) string { return shorten2(s, 400) }

func shorten2(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + " ..."
}

type Tool struct {
	Name    string `json:"name"`
	Found   bool   `json:"found"`
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
}

// Doctor reports which compilers are installed.
func Doctor() []Tool {
	var out []Tool
	for _, name := range []string{"javac", "java", "g++"} {
		t := Tool{Name: name}
		if p, err := exec.LookPath(name); err == nil {
			t.Found, t.Path = true, p
			args := []string{"--version"}
			if name == "javac" || name == "java" {
				args = []string{"-version"}
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			o, e, _ := execute(ctx, "", nil, name, args...)
			cancel()
			t.Version = firstLine(clean(o + e))
		}
		out = append(out, t)
	}
	return out
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
