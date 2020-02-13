package handlers

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/sessions"
	"github.com/go-redis/redis"
)

func TestGET(t *testing.T) {
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
	ctx := &ContextHandler{
		SigningKey:   signingKey,
		SessionStore: sessionStore,
		UserStore:    userStore,
	}
	mux := http.NewServeMux()
	wrappedMuxed := &CorsMW{Handler: mux}

}

func TestPOST(t *testing.T) {

}

func TestPATCH(t *testing.T) {

}

func TestDELETE(t *testing.T) {

}
