# AI Token Stats

Local token stats for a portfolio page. It reads the logs your AI coding tools already write and serves them over HTTP, so the page can draw its own card.

The server stays on `127.0.0.1`. It does not install hooks, and it does not upload the logs. One binary covers macOS, Linux, and Windows.

[![ci](https://github.com/zorexsh/ai-token-stats/actions/workflows/ci.yml/badge.svg)](https://github.com/zorexsh/ai-token-stats/actions/workflows/ci.yml)
[![license](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![go](https://img.shields.io/github/go-mod/go-version/zorexsh/ai-token-stats)](go.mod)

## Install

Download the archive for your machine from the [releases](https://github.com/zorexsh/ai-token-stats/releases), or build it. Go 1.22 or newer is enough. No C compiler.

```
go install github.com/zorexsh/ai-token-stats/cmd/ai-token-stats@latest
```

### macOS and Linux, Homebrew

The formula is `packaging/homebrew/ai-token-stats.rb`.

```
brew install --formula packaging/homebrew/ai-token-stats.rb
brew services start ai-token-stats
```

`brew services` starts it at login and restarts it if it exits.

### Linux, systemd

```
install -m 755 ai-token-stats /usr/local/bin/ai-token-stats
cp packaging/systemd/ai-token-stats.service ~/.config/systemd/user/
systemctl --user enable --now ai-token-stats
```

### Windows

The release archive ends in `.exe`. To start it at logon:

```
powershell -ExecutionPolicy Bypass -File packaging/windows/ai-token-stats.ps1 -Binary C:\Tools\ai-token-stats.exe
```

That registers a scheduled task named `ai-token-stats`. Remove it with `Unregister-ScheduledTask -TaskName ai-token-stats`.

## Run

```
ai-token-stats -port 7679
```

```
curl localhost:7679/stats?range=7d
```

`range` is `24h`, `7d`, `30d`, or `all`. The default is `7d`.

```json
{
  "totals": { "input": 0, "output": 0, "cache": 0, "tokens": 0, "usd": 0, "priced": true },
  "tools": [{ "tool": "claude", "totals": { "tokens": 0, "usd": 0, "priced": true } }],
  "models": [{ "model": "claude-opus-4-6", "totals": { "tokens": 0, "usd": 0, "priced": true } }],
  "days": [{ "day": "2026-11-10", "totals": { "tokens": 0, "usd": 0, "priced": true } }]
}
```

`usd` is an estimate from the published per-million rate for that model. Cache reads are billed at a tenth of input. `priced` is false when a model in the bucket has no published rate, and then `usd` is a lower bound, not a total. A route whose name ends in `-free` is priced at zero. The rates are in `internal/pricing.go`, with the date of the table in the comment above it.

`/health` returns `{"status":"ok"}`. Responses include `Access-Control-Allow-Origin: *`, so a page on another origin can fetch them.

## On a page

```html
<pre id="stats"></pre>
<script>
  fetch("http://127.0.0.1:7679/stats?range=30d")
    .then((r) => r.json())
    .then((d) => {
      document.getElementById("stats").textContent =
        d.totals.tokens.toLocaleString() + " tokens, $" + d.totals.usd.toFixed(2);
    });
</script>
```

## Logs

A tool with no folder is skipped. A file that does not match the expected shape is skipped too. Only token counts are kept. Prompt text never goes into the response.

| tool | macOS and Linux | Windows |
| --- | --- | --- |
| claude | `~/.claude/projects` | `%USERPROFILE%\.claude\projects` |
| codex | `~/.codex/sessions` | `%USERPROFILE%\.codex\sessions` |
| gemini | `~/.gemini` | `%USERPROFILE%\.gemini` |
| opencode | `$XDG_DATA_HOME/opencode` or `~/.local/share/opencode` | `%USERPROFILE%\.local\share\opencode` |
| cursor | `~/.cursor/projects`, plus `~/Library/Application Support/Cursor` on macOS and `$XDG_CONFIG_HOME/Cursor` on Linux | `%USERPROFILE%\.cursor\projects` and `%APPDATA%\Cursor` |
| grok | `~/.grok` | `%USERPROFILE%\.grok` |
| openclaw | `~/.openclaw` | `%USERPROFILE%\.openclaw` |
| cline | `~/.cline/data/sessions` | `%USERPROFILE%\.cline\data\sessions` |
| codebuddy | `~/.codebuddy` | `%USERPROFILE%\.codebuddy` |
| kimi | `~/.kimi` | `%USERPROFILE%\.kimi` |
| pi | `~/.pi` | `%USERPROFILE%\.pi` |

OpenCode follows `XDG_DATA_HOME` on every OS, including Windows. Cursor is also scanned in the per-OS application directory, so both locations are read.

## Build

```
make build
make test
make cross
```

`make cross` writes `dist/ai-token-stats-<os>-<arch>` for darwin, linux, and windows, arm64 and amd64. `CGO_ENABLED=0` is set because the SQLite driver is pure Go and does not need gcc.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). A new tool needs a reader and a row in the table above. A new model needs a row in `internal/pricing.go`.
