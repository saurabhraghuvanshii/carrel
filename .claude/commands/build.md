---
description: Implement a roadmap task from .plan/tasks (usage: /build 02)
argument-hint: <task number, for example 02>
---

Build task $ARGUMENTS.

1. Read `CLAUDE.md`, `.plan/START.md`, `.plan/PROGRESS.md`, and the task file `.plan/tasks/$ARGUMENTS-*.md`. Read any other `.plan` file the task points to.
2. If the task depends on another task that is not `done` in `.plan/PROGRESS.md`, or needs an answer from the owner, say so and stop.
3. Set `.plan/PROGRESS.md` to `doing` for this task. Write a short plan (files you will touch, in order) and carry on without waiting, unless something in the task is unclear or conflicts with the code. In that case ask one clear question first.
4. Implement every step in the task file. Follow the rules in `CLAUDE.md`. Write the tests the task asks for.
5. Run `gofmt -l .`, `go vet ./...` and `go test ./...`. Fix failures. Where the task describes an end-to-end check, run it with the real binary using `DSA_HOME` set to a temporary folder.
6. Go through the task's Acceptance list one by one and say which items you verified and how. Say plainly what you could not verify (for example a browser check).
7. Set the task to `done` in `.plan/PROGRESS.md`, add a short note, and update `.plan/ARCHITECTURE.md` or `.plan/DESIGN.md` if the work changed what they describe.
8. Show a short summary of what changed. Do not commit. Stop and wait for review.
