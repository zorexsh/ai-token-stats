# Contributing

Pull requests are welcome.

## Adding a tool

1. Find where the tool writes its usage. It is usually a `*.jsonl` session log or a SQLite file under the home directory.
2. Add a `source` entry in `internal/scan.go` with the folder and a reader.
3. If the file is JSONL and the usage sits in a `usage` object, `genericRecord` already handles it. Write a new reader only when the shape differs.
4. Add the tool to the table in `README.md`.

A reader returns token counts only. Never print or store the prompt or the response text.

## Checks

```
go test ./...
go build ./cmd/usage
```

Both should pass before you open the pull request.
