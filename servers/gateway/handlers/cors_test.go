package handlers

import (
	"testing"
)

var expectedHeaders = []string{
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Headers",
	"Access-Control-Expose-Headers",
	"Access-Control-Max-Age",
}

func TestGET(t *testing.T) {
	// signingKey := os.Getenv("SESSIONKEY")
	// redisAddr := os.Getenv("REDISADDR")
	// dsn := os.Getenv("DSN")

	// sessionStore := sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: redisAddr}), time.Hour)
	// if sessionStore == nil {
	// 	log.Fatal("Cannot connect to session store")
	// }
	// userStore, err := users.NewSQLStore(dsn)
	// if err != nil {
	// 	log.Fatal("Cannot connect to users database, reason: ", err)
	// }
	// ctx := &ContextHandler{
	// 	SigningKey:   signingKey,
	// 	SessionStore: sessionStore,
	// 	UserStore:    userStore,
	// }
	// mux := http.NewServeMux()
	// wrappedMuxed := &CorsMW{Handler: mux}

	// req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	// w := httptest.NewRecorder()

}

func TestPOST(t *testing.T) {

}

func TestPATCH(t *testing.T) {

}

func TestDELETE(t *testing.T) {

}
