package internal

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// A source is one tool and the function that reads a file from its logs.
type source struct {
	tool string
	read func(path, tool string, out *[]Record)
}

func sources() []source {
	return []source{
		{"claude", readJSONL(claudeRecord)},
		{"codex", readJSONL(codexRecord)},
		{"gemini", readJSONL(geminiRecord)},
		{"opencode", readSQLiteMessages},
		{"cursor", readJSONL(cursorRecord)},
		{"grok", readJSONL(genericRecord)},
		{"openclaw", readJSONL(genericRecord)},
		{"cline", readJSONL(genericRecord)},
		{"codebuddy", readJSONL(genericRecord)},
		{"kimi", readJSONL(genericRecord)},
		{"pi", readJSONL(genericRecord)},
	}
}

// Scan reads every supported tool's logs for the current machine.
func Scan() []Record {
	return scanLayout(hostLayout())
}

func scanLayout(l layout) []Record {
	var records []Record
	for _, s := range sources() {
		for _, dir := range l.dirs(s.tool) {
			filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() || !logFile(path) {
					return nil
				}
				s.read(path, s.tool, &records)
				return nil
			})
		}
	}
	return records
}

func logFile(path string) bool {
	name := strings.ToLower(filepath.Base(path))
	switch {
	case strings.HasSuffix(name, ".jsonl"):
		return true
	case strings.HasSuffix(name, ".db"):
		return true
	case strings.HasSuffix(name, ".sqlite"), strings.HasSuffix(name, ".sqlite3"):
		return true
	case strings.HasSuffix(name, ".vscdb"):
		return true
	default:
		return false
	}
}

func readJSONL(parse func(map[string]any, string) (Record, bool)) func(string, string, *[]Record) {
	return func(path, tool string, out *[]Record) {
		f, err := os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
		for sc.Scan() {
			var obj map[string]any
			if json.Unmarshal(sc.Bytes(), &obj) != nil {
				continue
			}
			rec, ok := parse(obj, tool)
			if ok && rec.Total() > 0 {
				*out = append(*out, rec)
			}
		}
	}
}

func claudeRecord(obj map[string]any, tool string) (Record, bool) {
	msg, _ := obj["message"].(map[string]any)
	if msg == nil {
		return Record{}, false
	}
	usage, _ := msg["usage"].(map[string]any)
	if usage == nil {
		return Record{}, false
	}
	return Record{
		Tool:   tool,
		Model:  str(msg["model"]),
		Hour:   hourOf(str(obj["timestamp"])),
		Input:  num(usage["input_tokens"]),
		Output: num(usage["output_tokens"]),
		Cache:  num(usage["cache_read_input_tokens"]) + num(usage["cache_creation_input_tokens"]),
	}, true
}

func codexRecord(obj map[string]any, tool string) (Record, bool) {
	payload, _ := obj["payload"].(map[string]any)
	if payload == nil {
		payload = obj
	}
	info, _ := payload["info"].(map[string]any)
	holder := info
	if holder == nil {
		holder = payload
	}
	usage, _ := holder["total_token_usage"].(map[string]any)
	if usage == nil {
		return genericRecord(obj, tool)
	}
	return Record{
		Tool:   tool,
		Model:  str(payload["model"]),
		Hour:   hourOf(str(obj["timestamp"])),
		Input:  num(usage["input_tokens"]),
		Output: num(usage["output_tokens"]),
		Cache:  num(usage["cached_input_tokens"]),
	}, true
}

func geminiRecord(obj map[string]any, tool string) (Record, bool) {
	tokens, _ := obj["tokens"].(map[string]any)
	if tokens == nil {
		return genericRecord(obj, tool)
	}
	return Record{
		Tool:   tool,
		Model:  str(obj["model"]),
		Hour:   hourOf(str(first(obj["timestamp"], obj["time"]))),
		Input:  num(tokens["input"]),
		Output: num(tokens["output"]),
		Cache:  num(tokens["cached"]),
	}, true
}

func cursorRecord(obj map[string]any, tool string) (Record, bool) {
	usage, _ := obj["usage"].(map[string]any)
	if usage == nil {
		return genericRecord(obj, tool)
	}
	return Record{
		Tool:   tool,
		Model:  str(obj["model"]),
		Hour:   hourOf(str(obj["timestamp"])),
		Input:  num(usage["input_tokens"]),
		Output: num(usage["output_tokens"]),
		Cache:  num(usage["cache_read_tokens"]),
	}, true
}

// genericRecord covers tools that log a flat usage object.
func genericRecord(obj map[string]any, tool string) (Record, bool) {
	usage, _ := obj["usage"].(map[string]any)
	if usage == nil {
		usage, _ = obj["token_usage"].(map[string]any)
	}
	if usage == nil {
		return Record{}, false
	}
	return Record{
		Tool:   tool,
		Model:  str(first(obj["model"], usage["model"])),
		Hour:   hourOf(str(first(obj["timestamp"], obj["time"]))),
		Input:  num(first(usage["input_tokens"], usage["input"])),
		Output: num(first(usage["output_tokens"], usage["output"])),
		Cache:  num(first(usage["cached_input_tokens"], usage["cache_read_tokens"])),
	}, true
}

func hourOf(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}
	}
	return parsed.UTC().Truncate(time.Hour)
}

func num(v any) int64 {
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	}
	return 0
}

func str(v any) string {
	s, _ := v.(string)
	if s == "" {
		return "unknown"
	}
	return s
}

func first(values ...any) any {
	for _, v := range values {
		if v != nil && v != "" {
			return v
		}
	}
	return nil
}
