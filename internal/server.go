package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	home string
}

func NewServer(home string) *Server { return &Server{home: home} }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /usage", s.usage)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	return mux
}

func (s *Server) usage(w http.ResponseWriter, r *http.Request) {
	from, to := window(r.URL.Query().Get("range"))
	records := Scan(s.home)
	writeJSON(w, Summarize(records, from, to))
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

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(v)
}
