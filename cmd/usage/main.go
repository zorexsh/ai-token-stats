package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/notzorexlol/usage/internal"
)

func main() {
	port := flag.Int("port", 7679, "port to listen on")
	flag.Parse()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	log.Printf("usage listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, internal.NewServer(home).Handler()))
}
