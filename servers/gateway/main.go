package main

import (
	"log"
	"net/http"
	"os"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/handlers"
)

//main is the main entry point for the server
func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":4000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary/", handlers.SummaryHandler)

	log.Printf("server is listening on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))

}
