package dao

import (
	"errors"
	"goTestProject/db"
	"time"
)

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Name     string
	Birthday time.Time
}

func (u *User) getTableName() string {
	return "user"
}

var mySqlDB = db.GetMySQLDb()

func (u *User) QueryUserByName(name string) User {
	var user User
	SelectByName(&user, name)
	return user
}

func (u *User) Add(name string, birthday time.Time) (User, error) {
	if name == "" {
		return User{}, errors.New("User name empty!")
	}
	user := u.QueryUserByName(name)
	if user.Id > 0 {
		return user, errors.New("User already exists!")
	}
	user.Name = name
	user.Birthday = birthday
	Insert(&user)
	return user, nil
}
