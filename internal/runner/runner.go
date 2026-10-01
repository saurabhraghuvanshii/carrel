// Package runner compiles and runs the user's code against a list of test
// cases using the compilers already installed on this computer.
//
// Test protocol: the driver reads "T" and then T cases from stdin and prints
// exactly one line of output per case, flushing after each one. A case that
// throws prints "ERROR ...".
//
// If the program dies part way, the case it was on is marked as crashed and
// the remaining cases run in a new process, at most maxRestarts times.
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
	"sync/atomic"
	"time"
)

const (
	maxRestarts = 5
	stdoutKeep  = 16 << 20
	stderrKeep  = 64 << 10
	outputLimit = 16 << 20 // a program that prints more than this on either stream is stopped
	memoryMB    = 1024     // C++ address space limit (ulimit -v), Unix only
)

const memoryMessage = "Your solution used too much memory."

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
	Note     string `json:"note,omitempty"` // why a case crashed or did not run
}

type Report struct {
	// Status is ok, compile_error, runtime_error, timeout, memory_limit,
	// tooling_missing or internal_error.
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
	RunTimeout     time.Duration // for all processes of one run together
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
			run:        []string{"java", "-Xss64m", "-Xmx256m", "-XX:+UseSerialGC", "-cp", "out", "Main"},
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
			run:        withMemoryLimit([]string{"." + string(filepath.Separator) + exe}, memoryMB),
		}, true
	}
	return spec{}, false
}

// Run compiles the code once and runs the cases, restarting after a crash.
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
	c := execute(cctx, dir, nil, sp.compile)
	if c.err != nil {
		if errors.Is(cctx.Err(), context.DeadlineExceeded) {
			return finish("timeout", fmt.Sprintf("Compiling took longer than %s, the compile time limit.", seconds(req.CompileTimeout)))
		}
		return finish("compile_error", clean(c.stdout+c.stderr))
	}

	rctx, rcancel := context.WithTimeout(ctx, req.RunTimeout)
	defer rcancel()

	n := len(req.Cases)
	got := make([]string, n)
	have := make([]bool, n)
	notes := make([]string, n)
	crashed := make([]bool, n)
	var firstCrash, crashStderr string
	memory, printedTooMuch, timedOut, restarted, gaveUp := false, false, false, false, false

	for first, restarts := 0, 0; first < n; restarts++ {
		var stdin bytes.Buffer
		fmt.Fprintf(&stdin, "%d\n", n-first)
		for _, c := range req.Cases[first:] {
			stdin.WriteString(c.Input)
			stdin.WriteString("\n")
		}
		r := execute(rctx, dir, &stdin, sp.run)

		var lines []string
		if s := strings.TrimRight(r.stdout, "\r\n"); s != "" {
			lines = strings.Split(s, "\n")
		}
		k := min(len(lines), n-first)
		for i := 0; i < k; i++ {
			got[first+i], have[first+i] = strings.TrimSpace(lines[i]), true
		}

		if r.overflow {
			printedTooMuch = true
			break
		}
		if r.outOfMemory {
			memory = true
			for i := first + k; i < n; i++ {
				have[i], got[i], notes[i] = true, "(not run)", "not run: an earlier case used too much memory"
			}
			break
		}
		if errors.Is(rctx.Err(), context.DeadlineExceeded) {
			timedOut = true
			break
		}
		if isMemory(r.stderr) {
			memory = true
		}
		if k == n-first {
			if r.err != nil && firstCrash == "" {
				firstCrash, crashStderr = "After the last case, your program "+describeExit(r.err, r.stderr)+".", r.stderr
			}
			break
		}

		// The program stopped before case first+k printed its line.
		at := first + k
		reason := describeExit(r.err, r.stderr)
		crashed[at], have[at], got[at], notes[at] = true, true, "(crashed)", reason
		if firstCrash == "" {
			firstCrash, crashStderr = fmt.Sprintf("%s: your program %s.", req.Cases[at].Label, reason), r.stderr
		}
		first = at + 1
		if first < n && restarts == maxRestarts {
			for i := first; i < n; i++ {
				have[i], got[i], notes[i] = true, "(not run)", fmt.Sprintf("not run: the program crashed %d times", maxRestarts+1)
			}
			gaveUp = true
			break
		}
		restarted = restarted || first < n
	}

	for i, c := range req.Cases {
		passed := have[i] && !crashed[i] && got[i] == strings.TrimSpace(c.Expected)
		r := Result{Kind: c.Kind, Label: c.Label, Passed: passed, Note: notes[i]}
		if !passed && r.Note == "" {
			r.Note = errorNote(got[i])
		}
		if isMemory(got[i]) {
			memory = true
		}
		if passed {
			rep.Passed++
		}
		if c.Kind == "example" || !passed {
			r.Input, r.Expected, r.Got = shorten(c.Input), shorten(c.Expected), shorten(got[i])
		}
		rep.Results = append(rep.Results, r)
	}

	switch {
	case timedOut:
		return finish("timeout", fmt.Sprintf("Your solution ran longer than %s, the time limit.", seconds(req.RunTimeout)))
	case printedTooMuch:
		return finish("runtime_error", fmt.Sprintf("Your program printed more than %d MB, so it was stopped. Remove prints inside loops.", outputLimit>>20))
	case memory:
		return finish("memory_limit", memoryMessage)
	case firstCrash != "":
		msg := firstCrash
		if restarted {
			msg += " The cases after it ran in a new process."
		}
		if gaveUp {
			msg += fmt.Sprintf(" After %d crashes the remaining cases were not run.", maxRestarts+1)
		}
		if s := clean(crashStderr); s != "" {
			msg += "\n" + s
		}
		return finish("runtime_error", msg)
	}
	return finish("ok", "")
}

func seconds(d time.Duration) string {
	return fmt.Sprintf("%g seconds", d.Seconds())
}

// isMemory spots out-of-memory errors in a result line or in stderr.
func isMemory(s string) bool {
	return strings.Contains(s, "OutOfMemoryError") || strings.Contains(s, "bad_alloc")
}

// errorNote explains an "ERROR ..." line printed by the driver.
func errorNote(got string) string {
	switch {
	case !strings.HasPrefix(got, "ERROR"):
		return ""
	case isMemory(got):
		return "used too much memory"
	case strings.Contains(got, "StackOverflowError"):
		return "stack overflow: the recursion went too deep"
	}
	return "your code threw an exception"
}

// describeExit turns how a process ended into a short phrase.
func describeExit(err error, stderr string) string {
	switch {
	case err == nil:
		return "stopped before printing a result (for example a call to exit)"
	case isMemory(stderr):
		return "used too much memory"
	case strings.Contains(stderr, "StackOverflowError"):
		return "hit a stack overflow: the recursion went too deep"
	}
	if msg, ok := describeSignal(err); ok {
		return msg
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return fmt.Sprintf("exited with code %d", ee.ExitCode())
	}
	return "stopped: " + err.Error()
}

type outcome struct {
	stdout, stderr string
	err            error
	overflow       bool // stopped for printing too much
	outOfMemory    bool // stopped after a case reported running out of memory
}

func execute(ctx context.Context, dir string, stdin *bytes.Buffer, argv []string) outcome {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var overflow, outOfMemory atomic.Bool
	stop := func() { overflow.Store(true); cancel() }
	memoryStop := func() { outOfMemory.Store(true); cancel() }

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.WaitDelay = 2 * time.Second
	if stdin != nil {
		cmd.Stdin = stdin
	}
	out := &limitedBuffer{keep: stdoutKeep, limit: outputLimit, onLimit: stop, onMemory: memoryStop}
	errb := &limitedBuffer{keep: stderrKeep, limit: outputLimit, onLimit: stop}
	cmd.Stdout, cmd.Stderr = out, errb
	err := cmd.Run()
	return outcome{out.buf.String(), errb.buf.String(), err, overflow.Load(), outOfMemory.Load()}
}

// limitedBuffer keeps the first keep bytes and calls onLimit once more than
// limit bytes have been written. If onMemory is set, it is called once when a
// result line reports running out of memory: the driver catches that error
// and would otherwise carry on into the next case and run out again.
// Each buffer is written by one goroutine.
type limitedBuffer struct {
	buf      bytes.Buffer
	keep     int
	limit    int
	total    int
	onLimit  func()
	onMemory func()
}

func (l *limitedBuffer) Write(p []byte) (int, error) {
	if room := l.keep - l.buf.Len(); room > 0 {
		l.buf.Write(p[:min(len(p), room)])
	}
	l.total += len(p)
	if l.total > l.limit && l.onLimit != nil {
		l.onLimit()
		l.onLimit = nil
	}
	if l.onMemory != nil && isMemory(string(p)) {
		l.onMemory()
		l.onMemory = nil
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
			args := []string{name, "--version"}
			if name == "javac" || name == "java" {
				args = []string{name, "-version"}
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			o := execute(ctx, "", nil, args)
			cancel()
			t.Version = firstLine(clean(o.stdout + o.stderr))
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
