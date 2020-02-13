package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/sessions"
)

// UsersHandler handles requests for the "users" resource
func (ctx *ContextHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		newUser := &users.NewUser{}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		err := json.Unmarshal(buf.Bytes(), newUser)
		if err != nil {
			http.Error(w, "Error unmarshalling new user from request: "+err.Error(), 500)
			return
		}
		// Make a new user and add to db
		user, err := newUser.ToUser()
		if err != nil {
			http.Error(w, "Error unmarshalling new user from request: "+err.Error(), 500)
			return
		}
		user, err = ctx.UserStore.Insert(user)
		if err != nil {
			http.Error(w, "Error inserting new user into the database: "+err.Error(), 500)
			return
		}
		// Begin session
		sessionState := SessionState{User: user, Time: time.Now()}
		_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sessionState, w)
		if err != nil {
			http.Error(w, "Error beginning session: "+err.Error(), 500)
			return
		}

		// Respond to client
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		http.Error(w, "Must be a post", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificUsersHandler handles requests for a specific user
func (ctx *ContextHandler) SpecificUsersHandler(w http.ResponseWriter, r *http.Request) {
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(w, "User not authenticated: "+err.Error(), http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "GET":
		var userID int64
		if base := path.Base(r.URL.Path); base == "me" {
			userID = sessionState.User.ID
		} else {
			userID, err = strconv.ParseInt(base, 10, 64)
			if err != nil {
				http.Error(w, "Cannot parse given user id", 400)
				return
			}
			if userID != sessionState.User.ID {
				http.Error(w, "Given ID does not match current authenticated ID", 400)
				return
			}
		}

		user, err := ctx.UserStore.GetByID(userID)
		if err != nil {
			http.Error(w, "UserID does not exist", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	case "PATCH":
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		base := path.Base(r.URL.Path)
		if base != "me" {
			// Assume base is the id and check if it matches
			incomingID, err := strconv.ParseInt(base, 10, 64)
			if err != nil {
				http.Error(w, "ID's do not match", http.StatusUnauthorized)
				return
			}
			if incomingID != sessionState.User.ID {
				http.Error(w, "ID's do not match", http.StatusUnauthorized)
				return
			}
		}

		// Assuming that this has no errors
		updates := &users.Updates{}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		err = json.Unmarshal(buf.Bytes(), updates)
		if err != nil {
			http.Error(w, "Error unmarshalling updates from request: "+err.Error(), 500)
			return
		}
		user, err := ctx.UserStore.Update(sessionState.User.ID, updates)
		if err != nil {
			http.Error(w, "Error updating user: "+err.Error(), 500)
			return
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	default:
		http.Error(w, "Must be a GET or PATCH request", http.StatusMethodNotAllowed)
		return
	}
}

// SessionsHandler handles requests for the "sessions" resource, and allows clients
// to begin a new session using an existing user's credentials.
func (ctx *ContextHandler) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		creds := &users.Credentials{}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		err := json.Unmarshal(buf.Bytes(), creds)
		if err != nil {
			http.Error(w, "Error unmarshalling creds from request: "+err.Error(), 500)
			return
		}
		user, err := ctx.UserStore.GetByEmail(creds.Email)
		if err != nil {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		err = user.Authenticate(creds.Password)
		if err != nil {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		// Begin a new session
		sessionState := SessionState{User: user, Time: time.Now()}
		_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sessionState, w)
		if err != nil {
			http.Error(w, "Error beginning session: "+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificSessionsHandler handles requests related to a specific authenticated session
func (ctx *ContextHandler) SpecificSessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		if strings.ToLower(path.Base(r.URL.Path)) != "mine" {
			http.Error(w, "Last path does not equal 'mine'", http.StatusForbidden)
			return
		}
		_, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("Signed Out"))
		return
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
}
