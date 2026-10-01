---
description: Add a new problem pack (usage: /add-problem <title> <easy|medium|hard> <patterns|real>)
---

Add a new problem pack for: $ARGUMENTS

Steps:

1. Read `CLAUDE.md` (problem rules) and the existing packs in `internal/problems/packs/` to copy their shape.
2. Choose an id in kebab-case and the right group and order for its sheet. Check `.plan/ROADMAP.md` for the planned learning order and the list of topics.
3. Write the statement from scratch in plain words. Do not copy from any website. Give it a real-world feel if the sheet is `real`.
4. Create `meta.json`, `statement.md` (with a `## Constraints` heading), `tests.json` (two visible examples and three or more edge cases), `starter.java`, `starter.cpp`, `driver.java`, `driver.cpp`. Follow the driver protocol in `CONTRIBUTING.md`.
5. Add a generator in `internal/testgen` and register it, with a test that checks its expected answers against a brute force.
6. Write the reference solutions `tools/refs/<id>/Solution.java` and `tools/refs/<id>/solution.cpp` (see `tools/refs/README.md`). Both are required.
7. Run `make check-packs` (or `go run ./cmd/packcheck <id>`). Every case must pass in both languages. Then submit one deliberately wrong solution through the running server and confirm it fails.
8. Run `gofmt -l .`, `go vet ./...` and `go test ./...`. `go test ./internal/problems` checks the pack structure: links in `buildsOn` and `leadsTo`, a `## Constraints` heading, at least 2 examples and 3 edge cases.

Report the id, the difficulty and how many cases each language passed.
