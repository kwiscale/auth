package auth

import (
	"os"
	"testing"

	"gopkg.in/kwiscale/orm.v0"
	_ "gopkg.in/kwiscale/orm.v0/dialects/sqlite"
)

func up() {
	orm.Initialize("sqlite", "test.db")
	Initialize()
}

func down() {
	os.Remove("test.db")
}

func TestCreateUser(t *testing.T) {
	up()
	defer down()

	u := User{}
	u.Login = "Foo"
	u.Password = "mypassword"
	err := u.Create()

	if err != nil {
		t.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	up()
	defer down()
	u := User{}
	u.Login = "Foo"
	u.Password = "mypassword"
	u.Create()

	u2 := User{}
	res := u2.Signin("Foo", "mypassword")
	if !res {
		t.Fatal("User should login")
	}

	u3 := User{}
	res = u3.Signin("Foo", "badpassword")

	if res {
		t.Fatal("User should NOT login")
	}
}

func TestUpdateUser(t *testing.T) {
	defer down()
	up()
	u := User{}
	u.Login = "Foo"
	u.Password = "mypassword"
	u.Create()

	test := User{}
	test.Signin("Foo", "mypassword")

	test.Password = "newpass"
	err := test.Save()
	if err != nil {
		t.Fatal("Password change error", err)
	}

	if !test.Signin("Foo", "newpass") {
		t.Fatal("User should be able to login with new password")
	}

	test.Login = "Foo2"
	err = test.Save()
	if err == nil {
		t.Fatal("User login cannot be changed")
	}

}

func TestRemoveUser(t *testing.T) {
	defer down()
	up()
	u := User{}
	u.Login = "Foo"
	u.Password = "mypassword"
	u.Create()

	err := u.Delete()
	if err != nil {
		t.Fatal("User deletion error", err)
	}
}

func TestInvalidSignup(t *testing.T) {

	u := User{}
	u.Password = "mypass"

	err := u.Create()
	if err == nil {
		t.Fatal("User shouldn't be able create without Login")
	}

	u.Login = "Foo"
	u.Password = ""
	err = u.Create()
	if err == nil {
		t.Fatal("User shouldn't be able create without password")
	}

	u.Login = "F"
	u.Password = "mypass"
	err = u.Create()
	if err == nil {
		t.Fatal("User shouldn't be able to create with a login length <", minlength)
	}

	u.Login = "Foo"
	u.Password = "m"
	err = u.Create()
	if err == nil {
		t.Fatal("User shouldn't be able to create with a password length too short")
	}
}
