package user

import (
	"errors"
)

var NotFoundError = errors.New("resource not found")

type UserRepository interface {
	Create(name, password string) error
	Update(id int64, name, password *string) error
	Delete(id int64) error
	Read(id *int64, username *string) (User, error)
}

type User struct {
	Username string
	Id       int64
}
