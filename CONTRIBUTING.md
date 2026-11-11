# Contributing

Pull requests are welcome. Open an issue first if the change is large, so the approach can be settled before the code is written.

## Adding a tool

1. Find where the tool writes its usage on macOS, Linux, and Windows. It is usually a `*.jsonl` session log or a SQLite file.
2. Add the directories to `internal/paths.go`. Tools that live under the home directory on every OS go in `dotDirs`. Tools that follow an OS convention, like Cursor, get their own case.
3. Add a `source` entry in `internal/scan.go`. If the file is JSONL and the usage sits in a `usage` object, `genericRecord` already handles it. Write a new reader only when the shape differs.
4. Add the tool to the table in `README.md`, with the path on each OS.
5. Add a fixture under `internal/paths_test.go` so the new directory is actually walked.
6. If the model has a published per-million rate, add it to `internal/pricing.go`. Leave it out if you would have to guess. A missing rate shows up as `priced: false`, which is the honest answer.

A reader returns token counts only. Never print or store the prompt or the response text.

## Checks

```
make test
make cross
```

`make cross` is the one that catches a build tag or a path separator that only fails on Windows. CI runs the tests on Ubuntu, macOS, and Windows.

## Commits

One concern per commit. Write the message as a sentence that says what changed and why, not a label.
