package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	every time.Duration
	read  func() []Record
	busy  atomic.Bool

	mu       sync.RWMutex
	records  []Record
	syncedAt time.Time
}

func NewServer(every time.Duration) *Server {
	s := &Server{every: every, read: (&Cache{}).Scan}
	s.refresh()
	go s.loop()
	return s
}

func (s *Server) refresh() {
	if !s.busy.CompareAndSwap(false, true) {
		return
	}
	defer s.busy.Store(false)

	start := time.Now()
	records := s.read()
	synced := time.Now().UTC()
	s.mu.Lock()
	s.records = records
	s.syncedAt = synced
	s.mu.Unlock()
	log.Printf("synced %d records in %s, next at %s", len(records), time.Since(start).Round(time.Millisecond), synced.Add(s.every).Format(time.RFC3339))
}

func (s *Server) loop() {
	tick := time.NewTicker(s.every)
	defer tick.Stop()
	for range tick.C {
		s.refresh()
	}
}

func (s *Server) snapshot() ([]Record, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Record, len(s.records))
	copy(out, s.records)
	return out, s.syncedAt
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stats", s.stats)
	mux.HandleFunc("GET /health", s.health)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Private-Network", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	from, to := window(r.URL.Query().Get("range"))
	records, synced := s.snapshot()
	body := Summarize(records, from, to)
	body.SyncedAt = synced.Format(time.RFC3339)
	body.NextSyncAt = synced.Add(s.every).Format(time.RFC3339)
	body.IntervalSeconds = int(s.every / time.Second)
	writeJSON(w, http.StatusOK, body)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	_, synced := s.snapshot()
	writeJSON(w, http.StatusOK, map[string]any{
		"status":           "ok",
		"synced_at":        synced.Format(time.RFC3339),
		"interval_seconds": int(s.every / time.Second),
		"next_sync_at":     synced.Add(s.every).Format(time.RFC3339),
	})
}

func window(name string) (time.Time, time.Time) {
	now := time.Now().UTC()
	to := now.Add(time.Hour)
	switch name {
	case "24h":
		return now.Add(-24 * time.Hour), to
	case "30d":
		return now.AddDate(0, 0, -30), to
	case "all":
		return time.Unix(0, 0), to
	default:
		return now.AddDate(0, 0, -7), to
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
