package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zorexsh/ai-token-stats/internal"
)

func main() {
	port := flag.Int("port", 7679, "port to listen on")
	every := flag.Duration("interval", 30*time.Second, "how often to reread the logs")
	flag.Parse()
	if *every <= 0 {
		log.Fatal("interval must be positive")
	}

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	log.Printf("ai-token-stats listening on http://%s, syncing every %s", addr, every.Truncate(time.Second))
	log.Fatal(http.ListenAndServe(addr, internal.NewServer(*every).Handler()))
}
