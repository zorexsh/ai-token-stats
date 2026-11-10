# usage

Reads the token logs your AI coding tools already write and serves the totals over HTTP, so a page can show them.

It does not install hooks and it does not phone home. It reads the files on disk and adds them up.

## Run

```
go install github.com/notzorexlol/usage/cmd/usage@latest
usage -port 7679
```

Then:

```
curl localhost:7679/usage?range=7d
```

`range` is `24h`, `7d`, `30d`, or `all`. The default is `7d`.

The response has `totals`, then the same numbers split by `tools`, `models`, and `days`.

```json
{
  "totals": { "input": 0, "output": 0, "cache": 0, "tokens": 0 },
  "tools": [{ "tool": "claude", "totals": { "tokens": 0 } }],
  "models": [{ "model": "claude-opus", "totals": { "tokens": 0 } }],
  "days": [{ "day": "2026-11-10", "totals": { "tokens": 0 } }]
}
```

`/health` returns `{"status":"ok"}`.

## What it reads

| tool | where |
| --- | --- |
| claude | `~/.claude/projects/**/*.jsonl` |
| codex | `~/.codex/sessions/**/*.jsonl` |
| gemini | `~/.gemini/**/*.jsonl` |
| opencode | `~/.local/share/opencode/opencode.db` |
| cursor | `~/.cursor/projects/**/*.jsonl` |
| grok | `~/.grok/**/*.jsonl` |
| openclaw | `~/.openclaw/**/*.jsonl` |
| cline | `~/.cline/data/sessions/**/*.jsonl` |
| codebuddy | `~/.codebuddy/**/*.jsonl` |
| kimi | `~/.kimi/**/*.jsonl` |
| pi | `~/.pi/**/*.jsonl` |

A tool with no folder is skipped. A file that does not match the expected shape is skipped too.

## Homebrew

```
brew install notzorexlol/tap/usage
brew services start usage
```

The formula lives in `packaging/homebrew/usage.rb`. It builds from source, so it works on macOS and Linux. On Windows, build with Go directly.

## Contributing

Open an issue or a pull request. A new tool needs a reader in `internal/scan.go` and a row in the table above. Keep the reader to token counts. Do not log prompt text.
