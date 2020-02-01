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
		addr = ":443"
	}
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	if (len(tlsKeyPath) == 0) || (len(tlsCertPath) == 0) {
		// If these environment variables are not set, write an error to standard out and exit the process with a non-zero code.
		log.Fatal("TLS Environment Variables Not Set")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary/", handlers.SummaryHandler)

	log.Printf("server is listening on port %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, mux))

}
