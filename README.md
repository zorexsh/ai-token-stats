# stats

[![ci](https://github.com/notzorexlol/stats/actions/workflows/ci.yml/badge.svg)](https://github.com/notzorexlol/stats/actions/workflows/ci.yml)
[![license](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![go](https://img.shields.io/github/go-mod/go-version/notzorexlol/stats)](go.mod)

Reads the token logs your AI coding tools already write and serves the stats over HTTP, so any page can show them.

It does not install hooks and it does not phone home. It reads the files on disk and adds them up. The same binary runs on macOS, Linux, and Windows.

## Install

Grab the archive for your machine from the [releases](https://github.com/notzorexlol/stats/releases), or build it yourself. Go 1.22 or newer is enough, and no C compiler is required.

```
go install github.com/notzorexlol/stats/cmd/stats@latest
```

### macOS and Linux, Homebrew

The formula is in `packaging/homebrew/stats.rb`. Point a tap at this repo, or install the file directly:

```
brew install --formula packaging/homebrew/stats.rb
brew services start stats
```

`brew services` keeps it running and restarts it after login.

### Linux, systemd

```
install -m 755 stats /usr/local/bin/stats
cp packaging/systemd/stats.service ~/.config/systemd/user/
systemctl --user enable --now stats
```

### Windows

The release archive ends in `.exe`. To start it at logon:

```
powershell -ExecutionPolicy Bypass -File packaging/windows/stats.ps1 -Binary C:\Tools\stats.exe
```

That registers a scheduled task. Remove it with `Unregister-ScheduledTask -TaskName stats`.

## Run

```
stats -port 7679
```

It binds to `127.0.0.1` only. Nothing outside the machine can reach it.

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

`usd` is an estimate from the published per-million rate for that model, with cache reads billed at a tenth of input. `priced` is false when a model in the bucket has no published rate, and `usd` is then a lower bound rather than a total. A route whose name ends in `-free` is treated as zero. Rates live in `internal/pricing.go` and are dated in the comment above the table.

`/health` returns `{"status":"ok"}`. The response carries `Access-Control-Allow-Origin: *`, so a page on any origin can fetch it.

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

## What it reads

A tool with no folder is skipped. A file that does not match the expected shape is skipped too. Token counts are the only thing kept. Prompt text is never read into the response.

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

OpenCode follows the XDG data directory on every operating system, Windows included. Set `XDG_DATA_HOME` and the database moves with it. Cursor also writes under the per-OS application directory, so both locations are scanned.

## Build from source

```
make build
make test
make cross
```

`make cross` writes `dist/stats-<os>-<arch>` for darwin, linux, and windows, on arm64 and amd64. `CGO_ENABLED=0` is set so the SQLite driver, which is pure Go, does not need gcc on the target.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). A new tool needs a reader and a row in the table above. A new model needs a row in `internal/pricing.go`.
