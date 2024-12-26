package user

import (
	"sync"

	"github.com/hashicorp/go-memdb"
)

func (db *MemDbUserRepository) Create(name, password string) error {
	txn := db.db.Txn(true)
	defer txn.Commit()

	db.idMutex.Lock()
	defer func() { db.id++; db.idMutex.Unlock() }()

	if err := txn.Insert("user", &User{
		Username: name,
		Id:       db.id,
		Password: password,
	}); err != nil {
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

// Read
// if id is supplied, find by id.
// else find by username and (optional) password
func (db *MemDbUserRepository) Read(id *int64, username *string, password *string) (User, error) {
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
	if password != nil && user.Password != *password {
		return User{}, PasswordError
	}
	return *user, nil
}

func NewMemDbUserRepository(db *memdb.MemDB) (*MemDbUserRepository, error) {
	return &MemDbUserRepository{
		db: db,
	}, nil
}

type MemDbUserRepository struct {
	db      *memdb.MemDB
	id      int64
	idMutex sync.Mutex
}
