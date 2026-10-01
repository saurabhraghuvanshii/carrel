package runner

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// A tiny problem for these tests: each case is a number x and solve(x)
// should return 2x. The drivers follow the real protocol.
const cppDriver = `#include <exception>
#include <iostream>
#include <string>
#include <vector>
using namespace std;

#include "solution.cpp"

int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        int x;
        cin >> x;
        try {
            out << solve(x) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
`

const javaDriver = `import java.io.*;

public class Main {
    public static void main(String[] args) throws IOException {
        PrintStream real = System.out;
        System.setOut(System.err);
        BufferedReader br = new BufferedReader(new InputStreamReader(System.in));
        int t = Integer.parseInt(br.readLine().trim());
        for (int c = 0; c < t; c++) {
            int x = Integer.parseInt(br.readLine().trim());
            String line;
            try {
                line = String.valueOf(new Solution().solve(x));
            } catch (Throwable e) {
                line = "ERROR " + e;
            }
            real.print(line + "\n");
            real.flush();
        }
    }
}
`

func need(t *testing.T, tools ...string) {
	t.Helper()
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
}

func run(t *testing.T, lang, code string, n int, limit time.Duration) Report {
	t.Helper()
	if lang == "cpp" {
		need(t, "g++")
	} else {
		need(t, "javac", "java")
	}
	driver := cppDriver
	if lang == "java" {
		driver = javaDriver
	}
	var cases []Case
	for x := 1; x <= n; x++ {
		cases = append(cases, Case{Kind: "edge", Label: "Case " + strconv.Itoa(x), Input: strconv.Itoa(x), Expected: strconv.Itoa(2 * x)})
	}
	rep := Run(context.Background(), Request{Lang: lang, Code: code, Driver: driver, Cases: cases,
		CompileTimeout: 60 * time.Second, RunTimeout: limit})
	if rep.Status == "compile_error" || rep.Status == "internal_error" {
		t.Fatalf("%s: %s", rep.Status, rep.Message)
	}
	return rep
}

func wantStatus(t *testing.T, rep Report, status, msgPart string) {
	t.Helper()
	if rep.Status != status || !strings.Contains(rep.Message, msgPart) {
		t.Fatalf("got %s %q, want %s containing %q", rep.Status, rep.Message, status, msgPart)
	}
}

func TestCorrectSolutionPasses(t *testing.T) {
	for lang, code := range map[string]string{
		"cpp":  `int solve(int x) { return 2 * x; }`,
		"java": `class Solution { int solve(int x) { return 2 * x; } }`,
	} {
		rep := run(t, lang, code, 10, 10*time.Second)
		if rep.Status != "ok" || rep.Passed != 10 {
			t.Fatalf("%s: %s %d/10 %s", lang, rep.Status, rep.Passed, rep.Message)
		}
	}
}

// checkCrashOnThree: cases 1 and 2 pass, 3 crashed with note, 4 to 10 pass.
func checkCrashOnThree(t *testing.T, rep Report, notePart string) {
	t.Helper()
	wantStatus(t, rep, "runtime_error", "new process")
	for i, r := range rep.Results {
		switch {
		case i == 2:
			if r.Passed || r.Got != "(crashed)" || !strings.Contains(r.Note, notePart) {
				t.Fatalf("case 3: %+v", r)
			}
		case !r.Passed:
			t.Fatalf("case %d should pass: %+v", i+1, r)
		}
	}
	if rep.Passed != 9 {
		t.Fatalf("passed %d, want 9", rep.Passed)
	}
}

func TestCppCrashOnCaseThreeRunsTheRest(t *testing.T) {
	code := `int solve(int x) {
    if (x == 3) {
        volatile int* p = nullptr;
        *p = 1;
    }
    return 2 * x;
}`
	checkCrashOnThree(t, run(t, "cpp", code, 10, 10*time.Second), "segmentation fault")
}

func TestJavaExitOnCaseThreeRunsTheRest(t *testing.T) {
	code := `class Solution {
    int solve(int x) {
        if (x == 3) System.exit(3);
        return 2 * x;
    }
}`
	checkCrashOnThree(t, run(t, "java", code, 10, 10*time.Second), "exited with code 3")
}

func TestGivesUpAfterFiveRestarts(t *testing.T) {
	code := `#include <cstdlib>
int solve(int x) { abort(); }`
	rep := run(t, "cpp", code, 10, 10*time.Second)
	wantStatus(t, rep, "runtime_error", "remaining cases were not run")
	for i, r := range rep.Results {
		want := "(crashed)"
		if i >= maxRestarts+1 {
			want = "(not run)"
		}
		if r.Got != want || r.Passed {
			t.Fatalf("case %d: %+v, want %s", i+1, r, want)
		}
	}
	if !strings.Contains(rep.Results[0].Note, "aborted") {
		t.Fatalf("abort not explained: %q", rep.Results[0].Note)
	}
}

func TestCppMemoryLimit(t *testing.T) {
	code := `#include <vector>
int solve(int x) {
    std::vector<std::vector<char>> hog;
    while (true) hog.emplace_back(1 << 20, 'x');
}`
	limit := 10 * time.Second
	rep := run(t, "cpp", code, 60, limit) // 60 hogs would overrun the time limit if each one ran
	wantStatus(t, rep, "memory_limit", memoryMessage)
	if time.Duration(rep.DurationMs)*time.Millisecond >= limit {
		t.Fatalf("took %d ms, not stopped by the memory limit", rep.DurationMs)
	}
}

func TestJavaMemoryLimit(t *testing.T) {
	code := `import java.util.*;
class Solution {
    int solve(int x) {
        List<int[]> hog = new ArrayList<>();
        while (true) hog.add(new int[1 << 20]);
    }
}`
	rep := run(t, "java", code, 30, 10*time.Second)
	wantStatus(t, rep, "memory_limit", memoryMessage)
	if rep.Results[0].Note != "used too much memory" {
		t.Fatalf("note: %q", rep.Results[0].Note)
	}
	if last := rep.Results[29]; last.Got != "(not run)" {
		t.Fatalf("cases after the first memory error should not run: %+v", last)
	}
}

func TestJavaStackOverflowIsExplained(t *testing.T) {
	code := `class Solution {
    int solve(int x) { return 1 + solve(x); }
}`
	rep := run(t, "java", code, 2, 10*time.Second)
	if rep.Passed != 0 || !strings.Contains(rep.Results[0].Note, "stack overflow") {
		t.Fatalf("%s %+v", rep.Status, rep.Results[0])
	}
}

func TestEndlessPrinterIsStopped(t *testing.T) {
	for name, code := range map[string]string{
		"cout (stderr)":   "#include <iostream>\nint solve(int x) { while (true) std::cout << \"debug line\\n\"; }",
		"printf (stdout)": "#include <cstdio>\nint solve(int x) { while (true) std::printf(\"debug line\\n\"); }",
	} {
		limit := 10 * time.Second
		rep := run(t, "cpp", code, 3, limit)
		if rep.Status != "runtime_error" || !strings.Contains(rep.Message, "printed more than") {
			t.Fatalf("%s: %s %q", name, rep.Status, rep.Message)
		}
		if time.Duration(rep.DurationMs)*time.Millisecond >= limit {
			t.Fatalf("%s: ran into the time limit instead of the output limit", name)
		}
	}
}

func TestInfiniteLoopTimesOut(t *testing.T) {
	code := `int solve(int x) { volatile int spin = 1; while (spin) {} return 0; }`
	rep := run(t, "cpp", code, 3, 2*time.Second)
	wantStatus(t, rep, "timeout", "2 seconds, the time limit")
}
