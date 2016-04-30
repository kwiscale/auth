package auth

import "gopkg.in/kwiscale/orm.v0"

func Initialize() {
	db := orm.Get()
	defer db.Close()
	db.AutoMigrate(&User{})
	db.Model(User{}).AddIndex("login_idx", "login")
}
