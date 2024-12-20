package user

import "github.com/hashicorp/go-memdb"

type UserRepository interface {
	Create(name, password string) error
	Update(id int64, name, password *string) error
	Delete(id int64) error
	Read(id int64) error
}

type MemDbUserRepository struct {
	db *memdb.MemDB
}

func (db *MemDbUserRepository) Update(id int64, name, password *string) error {
	//TODO implement me
	panic("implement me")
}

func (db *MemDbUserRepository) Delete(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (db *MemDbUserRepository) Read(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (db *MemDbUserRepository) Create(name, password string) error {
	return nil
}

func NewMemDbUserRepository(db *memdb.MemDB) (*MemDbUserRepository, error) {
	return &MemDbUserRepository{
		db: db,
	}, nil
}
