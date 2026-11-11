package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDirsFollowTheOS(t *testing.T) {
	home := filepath.Join(string(os.PathSeparator), "Users", "ada")
	l := layout{
		home:    home,
		appData: filepath.Join(home, "AppData", "Roaming"),
		config:  filepath.Join(home, ".config"),
		data:    filepath.Join(home, ".local", "share"),
	}

	l.goos = "windows"
	if got := l.dirs("cursor"); !hasSuffix(got, filepath.Join("AppData", "Roaming", "Cursor", "User", "workspaceStorage")) {
		t.Fatalf("windows cursor = %v", got)
	}

	l.goos = "linux"
	if got := l.dirs("cursor"); !hasSuffix(got, filepath.Join(".config", "Cursor", "User", "workspaceStorage")) {
		t.Fatalf("linux cursor = %v", got)
	}

	l.goos = "darwin"
	if got := l.dirs("cursor"); !hasSuffix(got, filepath.Join("Library", "Application Support", "Cursor", "User", "workspaceStorage")) {
		t.Fatalf("darwin cursor = %v", got)
	}

	if got := l.dirs("opencode"); !hasSuffix(got, filepath.Join(".local", "share", "opencode")) {
		t.Fatalf("opencode = %v", got)
	}
	if got := l.dirs("claude"); !hasSuffix(got, filepath.Join(".claude", "projects")) {
		t.Fatalf("claude = %v", got)
	}
}

func TestScanLayoutReadsEveryOS(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	l := layout{
		home:    home,
		appData: filepath.Join(home, "AppData", "Roaming"),
		config:  filepath.Join(home, ".config"),
		data:    filepath.Join(home, ".local", "share"),
		goos:    "windows",
	}

	writeJSONL(t, filepath.Join(home, ".claude", "projects", "sess", "a.jsonl"),
		`{"timestamp":"2026-11-10T03:12:00Z","message":{"model":"opus","usage":{"input_tokens":10,"output_tokens":4}}}`)
	writeJSONL(t, filepath.Join(l.appData, "Cursor", "User", "workspaceStorage", "ws", "chat.jsonl"),
		`{"timestamp":"2026-11-10T04:00:00Z","model":"cursor-small","usage":{"input_tokens":7,"output_tokens":1}}`)

	got := scanLayout(l)
	if len(got) != 2 {
		t.Fatalf("records = %+v", got)
	}
	if got[0].Tool == got[1].Tool {
		t.Fatalf("both records came from %s", got[0].Tool)
	}
}

func TestSQLiteDSNEscapesWindowsPaths(t *testing.T) {
	got := sqliteDSN(filepath.Join("C:", "Users", "ada", "opencode.db"))
	if strings.Contains(got, " ") || !strings.Contains(got, "mode=ro") || !strings.HasPrefix(got, "file:") {
		t.Fatalf("dsn = %s", got)
	}
	if strings.Contains(got, "\\") {
		t.Fatalf("dsn kept a backslash: %s", got)
	}
}

func writeJSONL(t *testing.T, path, line string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(line+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func hasSuffix(paths []string, suffix string) bool {
	for _, p := range paths {
		if len(p) >= len(suffix) && p[len(p)-len(suffix):] == suffix {
			return true
		}
	}
	return false
}
