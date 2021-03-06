package users

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		TestName     string
		Email        string
		Password     string
		PasswordConf string
		UserName     string
		FirstName    string
		LastName     string
		ExpectError  bool
	}{
		{
			"Working Validation",
			"alextan@gmail.com",
			"password",
			"password",
			"bigman69",
			"Alex",
			"Tan",
			false,
		},
		{
			"Invalid Email Address",
			"alextgmail.com",
			"password",
			"password",
			"bigman69",
			"Alex",
			"Tan",
			true,
		},
		{
			"Password too short",
			"alext@gmail.com",
			"pass",
			"pass",
			"bigman69",
			"Alex",
			"Tan",
			true,
		},
		{
			"Passwords do not match",
			"alext@gmail.com",
			"password",
			"passward",
			"bigman69",
			"Alex",
			"Tan",
			true,
		},
		{
			"Username cannot be empty",
			"alext@gmail.com",
			"password",
			"password",
			"",
			"Alex",
			"Tan",
			true,
		},
		{
			"Username cannot have spaces",
			"alext@gmail.com",
			"password",
			"password",
			"big man",
			"Alex",
			"Tan",
			true,
		},
	}

	for _, c := range cases {
		newUser := NewUser{
			Email:        c.Email,
			Password:     c.Password,
			PasswordConf: c.PasswordConf,
			UserName:     c.UserName,
			FirstName:    c.FirstName,
			LastName:     c.LastName,
		}
		err := newUser.Validate()
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error validating user: %v\n", c.TestName, err)
		}
		if err == nil && c.ExpectError {
			t.Errorf("case %s: expected error but didn't get one\n", c.TestName)
		}
	}
}

func TestToUser(t *testing.T) {
	cases := []struct {
		TestName     string
		Email        string
		Username     string
		Password     string
		PasswordConf string
		ExpectError  bool
	}{
		{
			"Correct email hashing for photoURL",
			"alextan785@gmail.com",
			"exampleusername",
			"password",
			"password",
			false,
		},
		{
			"Correct password hashing",
			"alextan785@gmail.com",
			"exampleusername",
			"password",
			"password",
			false,
		},
	}

	for _, c := range cases {
		newUser := NewUser{
			Email:        c.Email,
			UserName:     c.Username,
			Password:     c.Password,
			PasswordConf: c.PasswordConf,
		}
		user, err := newUser.ToUser()
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error validating user: %v\n", c.TestName, err)
		}
		// Check if c.Email is the same as photoURL after hashing
		photoURL := user.PhotoURL
		cleanEmail := strings.ToLower(strings.TrimSpace(c.Email))
		hasher := md5.New()
		hasher.Write([]byte(cleanEmail))
		testURL := gravatarBasePhotoURL + hex.EncodeToString(hasher.Sum(nil))
		if photoURL != testURL {
			t.Errorf("case %s: expected %s but got %s", c.TestName, testURL, photoURL)
		}

		// Check if the hashed password is correct
		err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(newUser.Password))
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error comparing password hash: %v\n", c.TestName, err)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		TestName      string
		FirstName     string
		LastName      string
		ExpectedValue string
	}{
		{
			"Valid Name",
			"Alex",
			"Tan",
			"Alex Tan",
		},
		{
			"Only First",
			"Alex",
			"",
			"Alex",
		},
		{
			"Only Last",
			"",
			"Tan",
			"Tan",
		},
		{
			"No Names",
			"",
			"",
			"",
		},
		{
			"Capitalize first letter",
			"alex",
			"tan",
			"Alex Tan",
		},
	}
	for _, c := range cases {
		user := User{
			FirstName: c.FirstName,
			LastName:  c.LastName,
		}
		fullName := user.FullName()
		if c.ExpectedValue != fullName {
			t.Errorf("case %s: expected %s but got %s\n", c.TestName, c.ExpectedValue, fullName)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	cases := []struct {
		TestName    string
		Password    string
		NewPassword string
		ExpectError bool
	}{
		{
			"Correct password",
			"password",
			"password",
			false,
		},
		{
			"Incorrect password",
			"password",
			"passward",
			true,
		},
		{
			"Zero length given password",
			"password",
			"",
			true,
		},
	}
	for _, c := range cases {
		user := User{}
		user.SetPassword(c.Password)
		err := user.Authenticate(c.NewPassword)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: password of %s and new password of %s caused an error\n", c.TestName, c.Password, c.NewPassword)
		}
		if err == nil && c.ExpectError {
			t.Errorf("case %s: expected error but didn't get one\n", c.TestName)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		TestName      string
		NewFirstName  string
		NewLastName   string
		ExpectedValue string
	}{
		{
			"Change Both First And Last",
			"John",
			"Smith",
			"John Smith",
		},
		{
			"Change Just First",
			"John",
			"",
			"John Tan",
		},
		{
			"Change Just Last",
			"",
			"Smith",
			"Alex Smith",
		},
		{
			"Change nothing",
			"",
			"",
			"Alex Tan",
		},
	}
	for _, c := range cases {
		user := User{
			FirstName: "Alex",
			LastName:  "Tan",
		}
		updates := &Updates{
			FirstName: c.NewFirstName,
			LastName:  c.NewLastName,
		}
		err := user.ApplyUpdates(updates)
		if err != nil {
			t.Errorf("case %s: didn't expect error but got one\n", c.TestName)
		}
		fullName := user.FullName()
		if fullName != c.ExpectedValue {
			t.Errorf("case %s: expected value %s but got %s\n", c.TestName, c.ExpectedValue, fullName)
		}
	}
}
