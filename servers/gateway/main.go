package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/handlers"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/sessions"
	"github.com/go-redis/redis"
)

// Director is a director
type Director func(r *http.Request)

// CustomDirector preprocesses the request for the microservice
// TODO: Check for current authenticated user?
func CustomDirector(targets []*url.URL) Director {
	var counter int32
	counter = 0
	return func(r *http.Request) {
		targ := targets[rand.Int()%len(targets)]
		atomic.AddInt32(&counter, 1)
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme
	}
}

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
	// TODO: Reverse proxy for summary and messages. How to convert this string to url?
	for _, port := range strings.Split(os.Getenv("SUMMARYADDR"), ",") {
		addr := "http:localhost:" + port
		summaryProxy := &httputil.ReverseProxy{Director: CustomDirector(addr)}
		mux.Handle("/v1/summary/", summaryProxy)
	}
	for _, port := range strings.Split(os.Getenv("MESSAGESADDR"), ",") {
		addr := "http:localhost:" + port
		messagesProxy := &httputil.ReverseProxy{Director: CustomDirector(addr)}
		mux.Handle("/v1/channels", messagesProxy)
		mux.Handle("/v1/messages", messagesProxy)
	}

	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUsersHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionsHandler)

	wrappedMuxed := &handlers.CorsMW{Handler: mux}

	log.Printf("server is listening on port %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMuxed))

}
