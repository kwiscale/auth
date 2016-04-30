package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/kwiscale/orm.v0"
)

const (
	minlength   = 3
	minpassword = 5
)

type User struct {
	orm.Model
	Login    string
	Password string
}

func (u *User) check() error {
	if len(u.Login) < minlength {
		return errors.New(fmt.Sprintf("Login size should be > %d", minlength))
	}
	if len(u.Password) < minpassword {
		return errors.New(fmt.Sprintf("Password length too short, should be length > %d", minpassword))
	}

	return nil
}

func (u *User) encryptPasswd() error {

	p, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err == nil {
		u.Password = string(p)
	}
	return err
}

func (u *User) Create() []error {
	err := u.check()
	if err != nil {
		return []error{err}
	}
	err = u.encryptPasswd()
	if err != nil {
		return []error{err}
	}
	d := orm.Get()
	defer d.Close()
	d.FirstOrCreate(u, map[string]string{"login": u.Login})
	return d.GetErrors()
}

func (u *User) Save() error {
	err := u.check()
	if err != nil {
		return err
	}

	db := orm.Get()
	defer db.Close()

	old := User{}
	db.Find(&old, map[string]uint{"id": u.ID})

	if u.Login != old.Login {
		return errors.New("Login cannot be changed")
	}

	if u.Password != "" && u.Password != old.Password {
		u.encryptPasswd()
	}
	db.Save(u)
	return nil
}

func (u *User) Delete() []error {
	d := orm.Get()
	defer d.Close()
	return d.Delete(u).GetErrors()
}

func (u *User) Signin(login, pass string) bool {
	db := orm.Get()
	defer db.Close()

	db.Find(u, map[string]string{"login": login})

	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)) == nil
}
