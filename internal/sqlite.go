package internal

import (
	"database/sql"
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// readSQLiteMessages reads OpenCode's message table. Tokens live in the JSON
// blob, not in columns, so each row is decoded on its own.
func readSQLiteMessages(path, tool string, out *[]Record) {
	if !isSQLiteFile(path) {
		return
	}
	db, err := sql.Open("sqlite", sqliteDSN(path))
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(`select time_created, data from message`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var created int64
		var data string
		if rows.Scan(&created, &data) != nil {
			continue
		}
		rec, ok := opencodeRecord(data, created, tool)
		if ok && rec.Total() > 0 {
			*out = append(*out, rec)
		}
	}
}

type opencodeMessage struct {
	ModelID    string `json:"modelID"`
	ProviderID string `json:"providerID"`
	Tokens     struct {
		Input  int64 `json:"input"`
		Output int64 `json:"output"`
		Cache  struct {
			Read int64 `json:"read"`
		} `json:"cache"`
	} `json:"tokens"`
}

func isSQLiteFile(path string) bool {
	name := strings.ToLower(path)
	return strings.HasSuffix(name, ".db") || strings.HasSuffix(name, ".sqlite") || strings.HasSuffix(name, ".sqlite3")
}

// sqliteDSN builds a read-only URI. A raw path breaks on Windows because the
// drive colon is read as a URI separator.
func sqliteDSN(path string) string {
	return "file:" + url.PathEscape(filepathToSlash(path)) + "?mode=ro"
}

func filepathToSlash(path string) string {
	return filepath.ToSlash(path)
}

func opencodeRecord(data string, created int64, tool string) (Record, bool) {
	var msg opencodeMessage
	if json.Unmarshal([]byte(data), &msg) != nil {
		return Record{}, false
	}
	model := msg.ModelID
	if msg.ProviderID != "" && msg.ProviderID != "unknown" {
		model = msg.ProviderID + "/" + msg.ModelID
	}
	return Record{
		Tool:   tool,
		Model:  model,
		Hour:   time.UnixMilli(created).UTC().Truncate(time.Hour),
		Input:  msg.Tokens.Input,
		Output: msg.Tokens.Output,
		Cache:  msg.Tokens.Cache.Read,
	}, true
}
