package user

import (
	"errors"
)

var NotFoundError = errors.New("resource not found")
var PasswordError = errors.New("password mismatch")

type UserRepository interface {
	Create(name, password string) error
	Update(id int64, name, password *string) error
	Delete(id int64) error
	Read(id *int64, username *string, password *string) (User, error)
}

type User struct {
	Username string
	Id       int64
	Password string
}
