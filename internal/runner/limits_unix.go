//go:build !windows

package runner

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

// withMemoryLimit runs argv under "ulimit -v". If the shell cannot set the
// limit (some macOS versions), the program still runs, just without it.
func withMemoryLimit(argv []string, mb int) []string {
	script := fmt.Sprintf(`ulimit -v %d 2>/dev/null; exec "$0" "$@"`, mb*1024)
	return append([]string{"sh", "-c", script}, argv...)
}

// describeSignal explains a process that was ended by a signal.
func describeSignal(err error) (string, bool) {
	var ee *exec.ExitError
	if !errors.As(err, &ee) {
		return "", false
	}
	ws, ok := ee.Sys().(syscall.WaitStatus)
	if !ok || !ws.Signaled() {
		return "", false
	}
	switch ws.Signal() {
	case syscall.SIGSEGV, syscall.SIGBUS:
		return "crashed with a segmentation fault: a bad memory access (an index out of range, a null pointer) or a stack overflow from very deep recursion", true
	case syscall.SIGFPE:
		return "crashed with an arithmetic error, usually a division by zero", true
	case syscall.SIGABRT:
		return "aborted: a failed assert, or an exception that is not a std::exception", true
	case syscall.SIGKILL:
		return "was killed (signal 9), often because the computer ran out of memory", true
	}
	return fmt.Sprintf("was stopped by signal %d (%s)", int(ws.Signal()), ws.Signal()), true
}
