package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"strings"

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
			w.Write([]byte("Error unmarshalling new user from request: " + err.Error()))
			return
		}
		// Begin session
		sessionState := SessionState{}
		sessionID, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sessionState, w)
		if err != nil {
			w.Write([]byte("Error beginning session: " + err.Error()))
			return
		}
		ctx.SessionStore.Save(sessionID, sessionState)
		// Make user and save into users db
		user, err := newUser.ToUser()
		if err != nil {
			w.Write([]byte("Error making new user: " + err.Error()))
			return
		}
		user, err = ctx.UserStore.Insert(user)
		if err != nil {
			w.Write([]byte("Error inserting new user into the database: " + err.Error()))
			return
		}
		// Respond to client
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

	} else {
		http.Error(w, "Must be a post", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificUsersHandler handles requests for a specific user
func (ctx *ContextHandler) SpecificUsersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Get userid
	switch r.Method {
	case "GET":
		// Get the user profile associated with the requested user ID from your store.
		// If no user is found with that ID, respond with an http.StatusNotFound (404) status code, and a suitable message.
		// Otherwise, respond to the client with:
		// a status code of http.StatusOK (200).
		// a response Content-Type header set to application/json to indicate that the response body is encoded as JSON
		// the users.User struct returned by your store in the response body, encoded as a JSON object.
		tempID := int64(1)
		user, err := ctx.UserStore.GetByID(tempID)
		if err != nil {
			http.Error(w, "UserID does not exist", http.StatusMethodNotAllowed)
			return
		}
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

	case "PATCH":
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		userID := path.Base(r.URL.Path)
		if userID != "me" {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		// TODO: match currently-authenticated user
		// TODO: Convert string to int64

		// Assuming that this has no errors
		updates := &users.Updates{}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		err := json.Unmarshal(buf.Bytes(), updates)
		if err != nil {
			w.Write([]byte("Error unmarshalling updates from request: " + err.Error()))
			return
		}
		// TODO: Update profile
		user, err := ctx.UserStore.Update(userID, updates)
		if err != nil {
			w.Write([]byte("Error updating user: " + err.Error()))
			return
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.Write([]byte(err.Error()))
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
			w.Write([]byte("Error unmarshalling creds from request: " + err.Error()))
			return
		}
		// Find and authenticate user
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
		sessionState := SessionState{}
		sessionID, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sessionState, w)
		if err != nil {
			w.Write([]byte("Error beginning session: " + err.Error()))
			return
		}
		ctx.SessionStore.Save(sessionID, sessionState)

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
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
		}
	// TODO: End session and respond with plain text message

	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
}
