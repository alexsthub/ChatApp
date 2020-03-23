package handlers

import (
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/indexes"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/sessions"
)

// ContextHandler stores the key used to sign and validate SessionIDs
// The session store to use when getting or saving session state
// The user store to use when finding or saving user profiles
type ContextHandler struct {
	SigningKey   string
	SessionStore *sessions.RedisStore
	UserStore    *users.MySQLStore
	UserTrie     *indexes.Trie
	Notifier     *Notifier
}
