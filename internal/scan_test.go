package internal

import (
	"testing"
	"time"
)

func TestSummarizeFiltersByWindow(t *testing.T) {
	hour := time.Date(2026, 11, 10, 3, 0, 0, 0, time.UTC)
	records := []Record{
		{Tool: "claude", Model: "opus", Hour: hour, Input: 100, Output: 20},
		{Tool: "opencode", Model: "gpt", Hour: hour.Add(-48 * time.Hour), Input: 5, Output: 5},
	}

	got := Summarize(records, hour.Add(-time.Hour), hour.Add(time.Hour))
	if got.Totals.Tokens != 120 {
		t.Fatalf("tokens = %d, want 120", got.Totals.Tokens)
	}
	if len(got.Tools) != 1 || got.Tools[0].Tool != "claude" {
		t.Fatalf("tools = %+v", got.Tools)
	}
}

func TestClaudeRecord(t *testing.T) {
	obj := map[string]any{
		"timestamp": "2026-11-10T03:12:00Z",
		"message": map[string]any{
			"model": "claude-opus",
			"usage": map[string]any{"input_tokens": 10.0, "output_tokens": 4.0, "cache_read_input_tokens": 2.0},
		},
	}
	rec, ok := claudeRecord(obj, "claude")
	if !ok || rec.Total() != 16 {
		t.Fatalf("record = %+v ok=%v", rec, ok)
	}
	if rec.Hour.Hour() != 3 {
		t.Fatalf("hour = %s", rec.Hour)
	}
}
