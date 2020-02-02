package users

import (
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	//- Email field must be a valid email address (hint: see mail.ParseAddress)
	//- Password must be at least 6 characters
	//- Password and PasswordConf must match
	//- UserName must be non-zero length and may not contain spaces
	//use fmt.Errorf() to generate appropriate error messages if
	//the new user doesn't pass one of the validation rules
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("Must be a valid email address")
	}
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password must be at least 6 characters")
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Password and password confirmation do not match")
	}
	if len(nu.UserName) == 0 || strings.Contains(" ", nu.UserName) {
		return fmt.Errorf("Username must not be empty and may not contain spaces")
	}
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	//return any validation errors that may occur.
	//if valid, create a new *User and set the fields
	//based on the field values in `nu`.
	//Leave the ID field as the zero-value; your Store
	//implementation will set that field to the DBMS-assigned
	//primary key value.
	//Set the PhotoURL field to the Gravatar PhotoURL
	//for the user's email address.
	//see https://en.gravatar.com/site/implement/hash/
	//and https://en.gravatar.com/site/implement/images/

	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:     nu.Email,
		UserName:  nu.UserName,
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		PhotoURL:  gravatarBasePhotoURL + nu.Email,
	}
	user.SetPassword(nu.Password)

	return user, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	if len(u.FirstName) != 0 && len(u.LastName) == 0 {
		return u.FirstName
	} else if len(u.FirstName) == 0 && len(u.LastName) != 0 {
		return u.LastName
	} else if len(u.FirstName) != 0 && len(u.LastName) != 0 {
		return u.FirstName + u.LastName
	} else {
		return ""
	}
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return err
	}
	u.PassHash = hash
	return nil

}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("Given password is empty")
	}
	if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password)); err != nil {
		return fmt.Errorf("password doesn't match stored hash")
	}
	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	// TODO: Errors What could go wrong?
	//set the fields of `u` to the values of the related
	//field in the `updates` struct
	if len(updates.FirstName) != 0 {
		u.FirstName = updates.FirstName
	}
	if len(updates.LastName) != 0 {
		u.LastName = updates.LastName
	}
	return nil
}
