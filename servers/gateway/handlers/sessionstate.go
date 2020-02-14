package handlers

import (
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/models/users"
)

//SessionState tracks the time at which this session began and the authenticated users who started the session
type SessionState struct {
	Time time.Time
	User *users.User
}
