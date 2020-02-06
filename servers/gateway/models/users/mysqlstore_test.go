package users

import (
	"testing"
)

// TODO: Do i need to insert rows before deleting / selecting?

func TestGetByID(t *testing.T) {
	cases := []struct {
		TestName    string
		ID          int
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
		store, err := NewSQLStore("users")
		if err != nil {
			t.Errorf("unexpected error getting database connection: %v", err)
		}
		_, err = store.GetByID(0)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error querying by id: %v", c.TestName, err)
		}
	}
}

func TestGetByEmail(t *testing.T) {

}

func TestGetByUsername(t *testing.T) {

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
			"alextan785@gmail.com",
			"Alex",
			"Tan",
			"alextan1000",
			"password",
			"password",
			false,
		},
		{
			"Null Values",
			"alextan785@gmail.com",
			"Alex",
			"Tan",
			"blahblahblah",
			"password",
			"password",
			true,
		},
	}

	for _, c := range cases {
		store, err := NewSQLStore("users")
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

		if c.TestName == "Null Values" {
			user.Email = ""
		}

		user, err = store.Insert(user)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error inserting: %v", c.TestName, err)
		}
	}
}

func TestUpdate(t *testing.T) {

}

func TestDelete(t *testing.T) {
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
		store, err := NewSQLStore("users")
		if err != nil {
			t.Errorf("unexpected error getting database connection: %v", err)
		}
		user := &User{
			Email:     "alextan785@gmail.com",
			PassHash:  []byte("AOSDNASODNASIDOASNNSOSAND"),
			UserName:  "MyUsername",
			FirstName: "Alex",
			LastName:  "Tan",
			PhotoURL:  "RandomPlaceholderPhotoURL",
		}
		user, err = store.Insert(user)
		if err != nil {
			t.Errorf("case %s: unexpected error inserting: %v", c.TestName, err)
		}

		err = store.Delete(c.ID)
		if err != nil && !c.ExpectError {
			t.Errorf("case %s: unexpected error deleting row: %v", c.TestName, err)
		}
	}
}
