package internal

import (
	"os"
	"path/filepath"
	"runtime"
)

// A layout is the set of directories one operating system uses for app data.
// Tests pass a fake one so the Windows and Linux layouts can be checked on
// any machine.
type layout struct {
	home    string
	appData string
	config  string
	data    string
	goos    string
}

func hostLayout() layout {
	home, _ := os.UserHomeDir()
	return layout{
		home:    home,
		appData: envOr("APPDATA", filepath.Join(home, "AppData", "Roaming")),
		config:  envOr("XDG_CONFIG_HOME", filepath.Join(home, ".config")),
		data:    envOr("XDG_DATA_HOME", filepath.Join(home, ".local", "share")),
		goos:    runtime.GOOS,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// dirs returns every directory a tool might keep logs in. Missing ones are
// skipped by the walk, so listing a fallback that does not exist is cheap.
func (l layout) dirs(tool string) []string {
	home := func(parts ...string) string {
		return filepath.Join(append([]string{l.home}, parts...)...)
	}
	switch tool {
	case "opencode":
		// OpenCode follows XDG on every platform, Windows included, and
		// writes a second database per release channel.
		return []string{filepath.Join(l.data, "opencode")}
	case "cursor":
		return []string{
			home(".cursor", "projects"),
			filepath.Join(l.cursorApp(), "User", "workspaceStorage"),
		}
	default:
		return []string{home(dotDirs[tool]...)}
	}
}

func (l layout) cursorApp() string {
	switch l.goos {
	case "windows":
		return filepath.Join(l.appData, "Cursor")
	case "darwin":
		return filepath.Join(l.home, "Library", "Application Support", "Cursor")
	default:
		return filepath.Join(l.config, "Cursor")
	}
}

// dotDirs are the tools that keep one folder under the home directory on
// every operating system. The slice is the path segments after home.
var dotDirs = map[string][]string{
	"claude":    {".claude", "projects"},
	"codex":     {".codex", "sessions"},
	"gemini":    {".gemini"},
	"grok":      {".grok"},
	"openclaw":  {".openclaw"},
	"cline":     {".cline", "data", "sessions"},
	"codebuddy": {".codebuddy"},
	"kimi":      {".kimi"},
	"pi":        {".pi"},
}
