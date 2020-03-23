package users

import (
	"database/sql"
	"fmt"
	"strings"

	// Sql driver
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-alexsthub/servers/gateway/indexes"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLStore represents a store for users
type MySQLStore struct {
	db *sql.DB
}

// NewSQLStore opens a connection and constructs a MySqlStore
func NewSQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	store := &MySQLStore{db: db}
	return store, nil
}

//GetByID returns the User with the given ID
func (store *MySQLStore) GetByID(id int64) (*User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	user := &User{}
	err := store.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PassHash,
		&user.UserName, &user.PhotoURL)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByEmail returns the User with the given email
func (store *MySQLStore) GetByEmail(email string) (*User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	user := &User{}
	err := store.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PassHash,
		&user.UserName, &user.PhotoURL)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetByUserName returns the User with the given Username
func (store *MySQLStore) GetByUserName(username string) (*User, error) {
	query := "SELECT * FROM users WHERE username = ?"
	user := &User{}
	err := store.db.QueryRow(query, username).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PassHash,
		&user.UserName, &user.PhotoURL)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (store *MySQLStore) Insert(user *User) (*User, error) {
	query := "INSERT INTO users (email, first_name, last_name, pass_hash, username, photo_url) VALUES (?, ?, ?, ?, ?, ?)"
	res, err := store.db.Exec(query, user.Email, user.FirstName, user.LastName, user.PassHash, user.UserName, user.PhotoURL)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (store *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	query := "UPDATE users SET first_name = ?, last_name = ? WHERE id = ?"
	_, err := store.db.Exec(query, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, err
	}
	user, err := store.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Delete deletes the user with the given ID
func (store *MySQLStore) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// trieUser is a struct
type trieUser struct {
	ID        int64
	FirstName string
	LastName  string
	UserName  string
}

// LoadUsersToTrie will load all users into a user trie and return it
func (store *MySQLStore) LoadUsersToTrie() (*indexes.Trie, error) {
	query := "SELECT id, first_name, last_name, username FROM users"
	rows, err := store.db.Query(query)
	if err != nil {
		return nil, err
	}
	userTrie := indexes.NewTrie()
	for rows.Next() {
		user := &trieUser{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName); err != nil {
			fmt.Printf("error scanning row: %v\n", err)
		}
		for _, word := range strings.Split(user.FirstName, " ") {
			userTrie.Add(strings.ToLower(word), user.ID)
		}
		for _, word := range strings.Split(user.LastName, " ") {
			userTrie.Add(strings.ToLower(word), user.ID)
		}
		for _, word := range strings.Split(user.UserName, " ") {
			userTrie.Add(strings.ToLower(word), user.ID)
		}
	}
	return userTrie, nil
}
