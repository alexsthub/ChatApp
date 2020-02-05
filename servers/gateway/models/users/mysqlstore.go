package users

// import (
// 	"database/sql"
// 	"fmt"
// 	"os"

// 	// Sql driver
// 	_ "github.com/go-sql-driver/mysql"
// )

// // MySQLStore represents a store for users
// type MySQLStore struct {
// 	db *sql.DB
// }

// // NewSQLStore opens a connection and constructs a MySqlStore
// func NewSQLStore(databaseName string) (*MySQLStore, error) {
// 	mySQLPassword := os.Getenv("MYSQL_ROOT_PASSWORD")
// 	if len(mySQLPassword) == 0 {
// 		return nil, fmt.Errorf("Password not found in environmental variables")
// 	}
// 	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/%s", mySQLPassword, databaseName)
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, fmt.Errorf("error opening database: %v", err)
// 	}

// 	//ensure that the database gets closed when we are done
// 	defer db.Close()
// 	if err := db.Ping(); err != nil {
// 		return nil, fmt.Errorf("error pinging database: %v", err)
// 	}
// 	store := &MySQLStore{db: db}
// 	return store, nil
// }

// //GetUserFromQuery takes a sql response and
// func GetUserFromQuery(res *sql.Rows) (*User, error) {
// 	defer res.Close()
// 	user := &User{}
// 	for res.Next() {
// 		if err := res.Scan(user.ID, user.Email, user.FirstName, user.LastName, user.PassHash,
// 			user.UserName, user.PhotoURL); err != nil {
// 			return nil, fmt.Errorf("error scanning row: %v", err)
// 		}
// 	}
// 	if err := res.Err(); err != nil {
// 		return nil, fmt.Errorf("error getting next row %v", err)
// 	}
// 	return user, nil
// }

// //GetByID returns the User with the given ID
// func (store *MySQLStore) GetByID(id int64) (*User, error) {
// 	query := "SELECT * FROM users WHERE id = ?"
// 	res, err := store.db.Query(query, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user, err := GetUserFromQuery(res)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// //GetByEmail returns the User with the given email
// func (store *MySQLStore) GetByEmail(email string) (*User, error) {
// 	query := "SELECT * FROM users WHERE email = '?'"
// 	res, err := store.db.Query(query, email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user, err := GetUserFromQuery(res)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// //GetByUserName returns the User with the given Username
// func (store *MySQLStore) GetByUserName(username string) (*User, error) {
// 	query := "SELECT * FROM users WHERE username = '?'"
// 	res, err := store.db.Query(query, username)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user, err := GetUserFromQuery(res)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// //Insert inserts the user into the database, and returns
// //the newly-inserted User, complete with the DBMS-assigned ID
// func (store *MySQLStore) Insert(user *User) (*User, error) {
// 	query := "INSERT INTO users (id, email, first_name, last_name, pass_hash, username, photo_url) VALUES (?, ?, ?, ?, ?, ?, ?)"
// 	_, err := store.db.Exec(query, user.ID, user.Email, user.FirstName, user.LastName, user.PassHash, user.UserName, user.PhotoURL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// //Update applies UserUpdates to the given user ID
// //and returns the newly-updated user
// func (store *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
// 	query := "UPDATE users SET first_name = '?', last_name = '?' WHERE id = ?"
// 	_, err := store.db.Exec(query, updates.FirstName, updates.LastName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user, err := store.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// //Delete deletes the user with the given ID
// func (store *MySQLStore) Delete(id int64) error {
// 	query := "DELETE FROM users WHERE id = ?"
// 	_, err := store.db.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
