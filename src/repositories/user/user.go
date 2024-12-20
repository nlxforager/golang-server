package user

import (
	"errors"
	"github.com/hashicorp/go-memdb"
	"sync"
)

var NotFoundError = errors.New("resource not found")

type UserRepository interface {
	Create(name, password string) error
	Update(id int64, name, password *string) error
	Delete(id int64) error
	Read(id *int64, username *string) (User, error)
}

type MemDbUserRepository struct {
	db      *memdb.MemDB
	id      int64
	idMutex sync.Mutex
}

type User struct {
	Username string
	Id       int64
}

func (db *MemDbUserRepository) Create(name, password string) error {
	txn := db.db.Txn(true)
	defer txn.Commit()

	db.idMutex.Lock()
	defer func() { db.id++; db.idMutex.Unlock() }()

	if err := txn.Insert("user", &User{Id: db.id, Username: name}); err != nil {
		txn.Abort()
		return err
	}

	return nil
}

func (db *MemDbUserRepository) Update(id int64, name, password *string) error {
	//TODO implement me
	panic("implement me")
}

func (db *MemDbUserRepository) Delete(id int64) error {

	txn := db.db.Txn(true)
	defer txn.Abort()
	err := txn.Delete("user", &User{Id: id})
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (db *MemDbUserRepository) Read(id *int64, username *string) (User, error) {
	//TODO implement me

	txn := db.db.Txn(false)
	defer txn.Abort()

	if id != nil {
		raw, err := txn.First("user", "id", *id)
		if err != nil {
			return User{}, err
		}
		user, ok := raw.(*User)
		if !ok {
			return User{}, NotFoundError
		}
		return *user, nil
	}

	raw, err := txn.First("user", "username", *username)
	if err != nil {
		return User{}, err
	}

	user := raw.(*User)
	return *user, nil
}

func NewMemDbUserRepository(db *memdb.MemDB) (*MemDbUserRepository, error) {
	return &MemDbUserRepository{
		db: db,
	}, nil
}
