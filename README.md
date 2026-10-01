# dsa

Practice data structures and algorithms in your browser. One small Go program, no account, nothing uploaded. Your solutions are plain files on your own computer.

```
dsa            # starts a local server and opens your browser
dsa doctor     # checks that Java and C++ compilers are installed
```

## Run it

You need Go 1.22 or newer, plus the compilers for the languages you want to practise in (`javac` and `java` for Java, `g++` for C++).

```
make build     # makes ./dsa
./dsa
make test      # unit tests
make check-packs   # runs reference solutions against every problem pack
make dist      # one binary per platform in dist/
```

Your data lives in `~/.dsa` (or `$DSA_HOME`):

```
~/.dsa/solutions/<problem>.java   your code, as plain files
~/.dsa/solutions/<problem>.cpp
~/.dsa/progress.json              tried / solved
~/.dsa/config.json                theme, accent colour, AI settings (owner-only permissions)
~/.dsa/packs/                     optional extra problem packs, same layout as below
```

## How it fits together

```
main.go                       start, bind 127.0.0.1 only, open the browser
cmd/packcheck/                checks every pack with the reference solutions
tools/refs/                   reference solutions (not shipped in the binary)
internal/server/              HTTP API and the embedded web UI (web/)
tools/vendor/                 builds the editor bundle and copies the fonts (Node, dev only)
internal/problems/            loads problem packs (embedded, plus ~/.dsa/packs)
internal/testgen/             seeded random test generators
internal/runner/              compile and run Java or C++ with a time limit
internal/store/               solutions and progress as plain files
internal/config/              settings file, API key included
internal/ai/                  Anthropic, OpenAI or Ollama, explain-only by default
```

Pressing Run or Submit:

1. The browser sends your code to the Go server.
2. The code is saved to `~/.dsa/solutions`.
3. The runner writes your code and the problem's driver to a temp folder, compiles once, and runs every case in one process.
4. Submit uses the examples, the fixed edge cases, and 50 fresh random cases from a new seed each time.
5. The results come back with the seed, so a failure can be reproduced.

## Problem packs

One folder per problem, in `internal/problems/packs/<id>/`:

| File | What it is |
| --- | --- |
| `meta.json` | title, difficulty (`easy`, `medium`, `hard`), sheet (`patterns` or `real`), group, order, tags, `buildsOn`, `leadsTo`, generator name |
| `statement.md` | the statement, in your own words |
| `tests.json` | visible `examples` and fixed `edge` cases |
| `starter.java`, `starter.cpp` | what the learner starts with |
| `driver.java`, `driver.cpp` | reads the test cases, calls the learner's code, prints one line per case |

The driver protocol is the same for every problem: stdin is `T` followed by `T` cases, stdout is exactly one line per case, and a case that throws prints `ERROR ...`. Anything the learner prints goes to stderr so it cannot break the results.

### Input and output formats

Every case's input uses one of these shapes, the same in every problem:

| Shape | Input |
| --- | --- |
| Array | `n`, then one line with the n values. An empty array is `0` and an empty line. |
| Linked list | Same as an array; the driver builds the nodes. Cycle problems add a line with the index the tail points to, or `-1`. |
| Binary tree | One line in level order, `null` for a missing child, trailing `null`s dropped, for example `5 3 8 null 4`. An empty tree is `null`. |
| Graph | `n m`, then m lines `u v` or `u v w`. |
| Grid | `rows cols`, then the rows as space-separated values. |
| Call sequence | `k`, then k lines `name arg ...`. The first call builds the object, for example `new 2`. |

Output is exactly one line per case and never an empty line. A list or array prints as `[a, b, c]`, an empty one as `[]`. A call sequence prints one result per call joined by spaces, with `null` for the constructor and for calls that return nothing. For a cache with room for 2 entries:

```
7
new 2
put 1 10
put 2 20
get 1
put 3 30
get 2
get 3
```

Expected line: `null null null 10 null -1 30`.

Copy-paste readers for each shape, in Java and C++, are in `internal/problems/drivers/README.md`.

`order` in `meta.json` is the group number times 100 plus the position in the group (`601` is the first linked-list problem), so a problem can be added without renumbering the others.

Random cases come from a Go generator registered in `internal/testgen`. A generator plants a known answer, or computes one with a reference solution, and returns the input and the expected line. Same seed, same cases.

To add a problem: copy a pack folder, change the files, write a generator, add reference solutions in `tools/refs/<id>/` (see `tools/refs/README.md`), then run `make test` and `make check-packs`.

`make test` checks the structure of every pack: all files present, a registered generator, a `## Constraints` heading, at least 2 examples and 3 edge cases, no trailing spaces in inputs, every `buildsOn` and `leadsTo` id exists, `leadsTo` matches `buildsOn`, and no cycles.

`make check-packs` runs the Java and C++ reference solutions through the real runner on the examples, the edge cases and 200 random cases from a fixed seed. Every case must pass in both languages, so a wrong expected value or a broken driver is caught before it ships. It needs `javac`, `java` and `g++`.

## Rebuilding the editor bundle

The editor (CodeMirror 6 with Java and C++) and the three fonts are committed under `internal/server/web/vendor/` and `internal/server/web/fonts/`, so the app works offline and Go never needs Node. To rebuild them you need Node:

```
cd tools/vendor
npm ci
npm run build
```

See `tools/vendor/README.md` for what the bundle exposes. The font licences are in `internal/server/web/fonts/LICENSES.md`.

## Safety

This program runs code on your computer, so the server is locked down:

- it listens on `127.0.0.1` only
- requests for any other Host are refused, which blocks DNS rebinding
- anything that changes state needs an `X-DSA` header that other websites cannot add
- problem ids and languages are validated before they touch a file path
- the API key is never sent back to the browser

Your own code is not sandboxed beyond a time limit, the same as running it yourself.

## Not done yet

- Memory limit for the code you run (a time limit exists)
- Function-style wrappers for more languages (Python, JavaScript, Go)
- `dsa pull` to download extra problem packs
- Export and import of all solutions as one zip
- More problems: the aim is about 100 patterns problems and a real-interview sheet
