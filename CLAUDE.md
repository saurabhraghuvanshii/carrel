# Carrel

A local DSA practice tool. One Go binary starts a server on `127.0.0.1`, opens the browser, runs the user's Java or C++ code with the compilers already on their computer, and saves solutions as plain files in `~/.carrel`. No account, no cloud, no Docker.

## Start of every session

Private planning notes live in `.plan/` (listed in `.gitignore`, never pushed). Read `.plan/START.md` first, then `.plan/ROADMAP.md`. The slash command `/resume` does this for you. If `.plan/` is missing, ask the owner for it before making product decisions.

## Commands

```
make build        # builds ./carrel
make test         # go test ./...
make vet          # go vet ./...
make check-packs  # reference solutions in tools/refs must pass every pack
make site         # builds the website into site/dist (make site-serve to preview)
./carrel doctor      # shows which compilers are installed
./carrel --no-open --port 7777
```

Run `gofmt -l .`, `go vet ./...` and `go test ./...` before saying a change is done. CARREL_HOME=<dir> moves the data folder, which keeps manual testing away from the real `~/.carrel`.

## Layout

```
main.go                  start-up, loopback listener, `doctor` subcommand
cmd/packcheck/           pack checker behind `make check-packs`
tools/refs/              reference solutions per problem, never embedded
internal/server/         HTTP API + embedded UI (web/index.html, style.css, app.js)
internal/problems/       loads packs from packs/ (embedded) and ~/.carrel/packs
internal/testgen/        seeded random generators, registered by name
internal/runner/         compile once, run all cases in one process, time limits
internal/store/          solutions and progress as plain files
internal/config/         ~/.carrel/config.json (owner-only; holds the AI key)
internal/ai/             Anthropic, OpenAI, Ollama; explain-only by default
site/                    public website source (built with `make site`)
tools/sitebuild/         Go tool that assembles site/src into site/dist
scripts/                 install.sh and install.ps1 served by the website
```

## Rules for this codebase

- Standard library only. Ask the owner before adding a dependency.
- The UI is plain HTML, CSS and JavaScript with no build step and no CDN. Everything must work offline.
- Build DOM with `textContent` or the `el()` helper, never `innerHTML` with user or problem text.
- Safety checks stay in place: loopback only, Host check, `X-Carrel` header on writes, validated ids and languages, the API key never returned to the browser. Do not weaken them.
- Every new feature gets a test. Server behaviour is tested with `httptest`, runner behaviour with real `javac`, `java` and `g++` where available.
- Keep files small and plain. No comments that only restate the code.

## Problems

- Statements are always written from scratch. Never copy text, examples or test data from LeetCode, GeeksforGeeks or any other site. Patterns and algorithms are free to use; wording is not.
- Do not claim a problem was asked by a named company. Say "pattern commonly seen in online assessments".
- Pack layout and driver protocol are in `CONTRIBUTING.md`. A pack must have both languages, a generator name that exists in `internal/testgen`, a `## Constraints` heading, at least 2 examples and at least 3 edge cases.
- `go test ./internal/problems` fails on an incomplete pack, and `make check-packs` fails when a reference solution does not pass every case. Keep both passing.

## Design rules

The UI follows two themes, Paper (warm off-white) and Ink (warm charcoal), with five accent colours and fixed Easy / Medium / Hard tag colours. No gradients, shadows, glass effects or emoji. No blue-on-dark "AI look". Sentence case everywhere. Details are in `.plan/DESIGN.md`.

## Commits

Small commits with a plain sentence in the subject. Never commit `.plan/`, `~/.carrel` contents, API keys or built binaries.
