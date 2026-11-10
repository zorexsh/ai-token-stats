package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/notzorexlol/stats/internal"
)

func main() {
	port := flag.Int("port", 7679, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	log.Printf("stats listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, internal.NewServer().Handler()))
}
