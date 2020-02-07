package handlers

import (
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

// ContextHandler stores the key used to sign and validate SessionIDs
// The session store to use when getting or saving session state
// The user store to use when finding or saving user profiles
type ContextHandler struct {
	SigningKey   string
	SessionStore *sessions.RedisStore
	UserStore    *users.MySQLStore
}
