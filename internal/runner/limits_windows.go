//go:build windows

package runner

// withMemoryLimit does nothing on Windows yet. A job object would be needed,
// so C++ solutions there only have the time limit. Java still has -Xmx.
func withMemoryLimit(argv []string, mb int) []string { return argv }

// describeSignal: Windows has no signals; exit codes are reported instead.
func describeSignal(err error) (string, bool) { return "", false }
