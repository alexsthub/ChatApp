package users

import (
	"testing"
)

func TestGetByID(t *testing.T) {
	cases := []struct {
		TestName    string
		ID          int64
		ExpectError bool
	}{
		{
			"Correct userid",
			1,
			false,
		},
		{
			"User ID does not exist",
			1000000000,
			true,
		},
	}
	for _, c := range cases {
		store, err := NewSQLStore("root:password@tcp(127.0.0.1:3306)/users")
		if err != nil {
			t.Errorf("unexpected error getting database connection: %v", err)
		}
		_, err = store.GetByID(c.ID)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error querying by id: %v", c.TestName, err)
		}
	}
}

func TestGetByEmail(t *testing.T) {
	cases := []struct {
		TestName    string
		Email       string
		ExpectError bool
	}{
		{
			"Correct userid",
			"alextan785@gmail.com",
			false,
		},
		{
			"User ID does not exist",
			"test@test.com",
			true,
		},
	}
	for _, c := range cases {
		store, err := NewSQLStore("root:password@tcp(127.0.0.1:3306)/users")
		if err != nil {
			t.Errorf("unexpected error getting database connection: %v", err)
		}
		_, err = store.GetByEmail(c.Email)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error querying by email: %v", c.TestName, err)
		}
	}
}

func TestGetByUsername(t *testing.T) {
	cases := []struct {
		TestName    string
		UserName    string
		ExpectError bool
	}{
		{
			"Correct username",
			"alextan785",
			false,
		},
		{
			"User ID does not exist",
			"doomedtofail",
			true,
		},
	}
	for _, c := range cases {
		store, err := NewSQLStore("root:password@tcp(127.0.0.1:3306)/users")
		if err != nil {
			t.Errorf("unexpected error getting database connection: %v", err)
		}
		_, err = store.GetByUserName(c.UserName)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error querying by username: %v", c.TestName, err)
		}
	}
}

func TestInsert(t *testing.T) {
	cases := []struct {
		TestName     string
		Email        string
		FirstName    string
		LastName     string
		Username     string
		Password     string
		PasswordConf string
		ExpectError  bool
	}{
		{
			"Working insert",
			"butt@gmail.com",
			"Alex",
			"Tan",
			"alextan1000",
			"password",
			"password",
			false,
		},
	}

	for _, c := range cases {
		store, err := NewSQLStore("root:password@tcp(127.0.0.1:3306)/users")
		if err != nil {
			t.Errorf("unexpected error getting database connection: %v", err)
		}
		newUser := NewUser{
			Email:        c.Email,
			UserName:     c.Username,
			Password:     c.Password,
			PasswordConf: c.PasswordConf,
			FirstName:    c.FirstName,
			LastName:     c.LastName,
		}
		user, err := newUser.ToUser()
		if err != nil {
			t.Errorf("case %s: problem making new user: %v", c.TestName, err)
		}

		user, err = store.Insert(user)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error inserting: %v", c.TestName, err)
		}
		if err == nil && c.ExpectError {
			t.Errorf("case %s: expected error inserting row: %v", c.TestName, err)
		}
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		TestName        string
		FirstNameUpdate string
		LastNameUpdate  string
	}{
		{
			"Update first and last",
			"Peter",
			"Long",
		},
		{
			"Update just first",
			"Jerry",
			"",
		},
		{
			"Update last",
			"",
			"Lee",
		},
		{
			"Update none",
			"",
			"",
		},
	}
	store, err := NewSQLStore("root:password@tcp(127.0.0.1:3306)/users")
	if err != nil {
		t.Errorf("unexpected error getting database connection: %v", err)
	}
	for _, c := range cases {
		update := &Updates{
			FirstName: c.FirstNameUpdate,
			LastName:  c.LastNameUpdate,
		}
		fakeUser := User{
			FirstName: c.FirstNameUpdate,
			LastName:  c.LastNameUpdate,
		}
		user, err := store.Update(2, update)
		if err != nil {
			t.Errorf("unexpected error updating row: %v", err)
		}
		if user.FullName() != fakeUser.FullName() {
			t.Errorf("case %s: expected %s but got %s", c.TestName, fakeUser.FullName(), user.FullName())
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		TestName    string
		ID          int64
		ExpectError bool
	}{
		{
			"Correct userid",
			0,
			false,
		},
		{
			"User ID does not exist",
			1000000000,
			false,
		},
	}

	for _, c := range cases {
		store, err := NewSQLStore("root:password@tcp(127.0.0.1:3306)/users")
		if err != nil {
			t.Errorf("unexpected error getting database connection: %v", err)
		}
		err = store.Delete(c.ID)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error deleting row: %v", c.TestName, err)
		}
		if err == nil && c.ExpectError {
			t.Errorf("case %s: expected error deleting row: %v", c.TestName, err)
		}
	}
}
