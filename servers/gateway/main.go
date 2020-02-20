package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/handlers"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/sessions"
	"github.com/go-redis/redis"
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
		log.Fatal("TLS Environment Variables Not Set")
	}

	signingKey := os.Getenv("SESSIONKEY")
	redisAddr := os.Getenv("REDISADDR")
	dsn := os.Getenv("DSN")

	sessionStore := sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: redisAddr}), time.Hour)
	if sessionStore == nil {
		log.Fatal("Cannot connect to session store")
	}
	userStore, err := users.NewSQLStore(dsn)
	if err != nil {
		log.Fatal("Cannot connect to users database, reason: ", err)
	}
	ctx := &handlers.ContextHandler{
		SigningKey:   signingKey,
		SessionStore: sessionStore,
		UserStore:    userStore,
	}

	// Load users into trie
	userTrie, err := ctx.UserStore.LoadUsersToTrie()
	if err != nil {
		log.Fatal("Error loading users to user trie")
	}
	ctx.UserTrie = userTrie

	mux := http.NewServeMux()
	// Reverse proxy for summary
	summaryProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "localhost:80"})
	mux.Handle("/v1/summary/", summaryProxy)

	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUsersHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionsHandler)

	wrappedMuxed := &handlers.CorsMW{Handler: mux}

	log.Printf("server is listening on port %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMuxed))

}
