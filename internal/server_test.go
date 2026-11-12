package internal

import (
	"testing"
	"time"
)

func TestRefreshKeepsTheSnapshot(t *testing.T) {
	hour := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	s := &Server{
		every: time.Hour,
		read: func() []Record {
			return []Record{{Tool: "claude", Model: "test", Hour: hour, Input: 4, Output: 1}}
		},
	}
	s.refresh()

	got, synced := s.snapshot()
	if len(got) != 1 || got[0].Input != 4 {
		t.Fatalf("snapshot = %+v", got)
	}
	if synced.IsZero() {
		t.Fatal("synced_at was not set")
	}
}
