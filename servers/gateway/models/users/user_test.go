package users

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

// func testValidate(t *testing.T) {
// 	cases := []struct {
// 		TestName     string
// 		Email        string
// 		Password     string
// 		PasswordConf string
// 		UserName     string
// 		FirstName    string
// 		LastName     string
// 		ExpectError  bool
// 	}{
// 		{
// 			"Working Validation",
// 			"alextan@gmail.com",
// 			"password",
// 			"password",
// 			"bigman69",
// 			"Alex",
// 			"Tan",
// 			false,
// 		},
// 		{
// 			"Invalid Email Address",
// 			"alextgmail.com",
// 			"password",
// 			"password",
// 			"bigman69",
// 			"Alex",
// 			"Tan",
// 			true,
// 		},
// 		{
// 			"Password too short",
// 			"alext@gmail.com",
// 			"pass",
// 			"pass",
// 			"bigman69",
// 			"Alex",
// 			"Tan",
// 			true,
// 		},
// 		{
// 			"Passwords do not match",
// 			"alext@gmail.com",
// 			"password",
// 			"passward",
// 			"bigman69",
// 			"Alex",
// 			"Tan",
// 			true,
// 		},
// 	}

// 	for _, c := range cases {
// 		newUser := NewUser{
// 			Email:        c.Email,
// 			Password:     c.Password,
// 			PasswordConf: c.PasswordConf,
// 			UserName:     c.UserName,
// 			FirstName:    c.FirstName,
// 			LastName:     c.LastName,
// 		}
// 		err := newUser.Validate()
// 		if err != nil && !c.ExpectError {
// 			t.Errorf("case %s: unexpected error validating user: %v\n", c.TestName, err)
// 		}
// 		if err == nil && c.ExpectError {
// 			t.Errorf("case %s: expected error but didn't get one\n", c.TestName)
// 		}
// 	}
// }

// // TODO
// func testToUser(t *testing.T) {

// }

// func testFullName(t *testing.T) {
// 	cases := []struct {
// 		TestName      string
// 		FirstName     string
// 		LastName      string
// 		ExpectedValue string
// 	}{
// 		{
// 			"Valid Name",
// 			"Alex",
// 			"Tan",
// 			"Alex Tan",
// 		},
// 		{
// 			"Only First",
// 			"Alex",
// 			"",
// 			"Alex",
// 		},
// 		{
// 			"Only Last",
// 			"",
// 			"Tan",
// 			"Tan",
// 		},
// 		{
// 			"No Names",
// 			"",
// 			"",
// 			"",
// 		},
// 	}
// 	for _, c := range cases {
// 		user := User{
// 			FirstName: c.FirstName,
// 			LastName:  c.LastName,
// 		}
// 		fullName := user.FullName()
// 		if c.ExpectedValue != fullName {
// 			t.Errorf("case %s: expected %s but got %s\n", c.TestName, c.ExpectedValue, fullName)
// 		}
// 	}
// }

// func testAuthenticate(t *testing.T) {
// 	cases := []struct {
// 		TestName    string
// 		Password    string
// 		NewPassword string
// 		ExpectError bool
// 	}{
// 		{
// 			"Correct password",
// 			"password",
// 			"password",
// 			false,
// 		},
// 		{
// 			"Incorrect password",
// 			"password",
// 			"passward",
// 			true,
// 		},
// 		{
// 			"Zero length given password",
// 			"password",
// 			"",
// 			true,
// 		},
// 	}
// 	for _, c := range cases {
// 		user := User{}
// 		user.SetPassword(c.Password)
// 		err := user.Authenticate(c.NewPassword)
// 		if err != nil && !c.ExpectError {
// 			t.Errorf("case %s: password of %s and new password of %s caused an error\n", c.TestName, c.Password, c.NewPassword)
// 		}
// 		if err == nil && c.ExpectError {
// 			t.Errorf("case %s: expected error but didn't get one\n", c.TestName)
// 		}
// 	}
// }

// // TODO
// func testApplyUpdates(t *testing.T) {

// }
