# Reference solutions

One folder per problem, named like its pack:

```
tools/refs/<problem-id>/Solution.java
tools/refs/<problem-id>/solution.cpp
```

These are correct solutions used only by the pack checker. They are not embedded in the binary and users never see them.

To add a reference:

1. Write `Solution.java` with `class Solution` and `solution.cpp` with the same free functions as the pack's starters.
2. Run `make check-packs` (or `go run ./cmd/packcheck <problem-id>` for one pack).

The checker runs both references through the real runner on the examples, the edge cases and 200 random cases from a fixed seed. Every case must pass in both languages. A missing reference counts as a failure.
