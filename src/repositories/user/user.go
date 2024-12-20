package user

type UserRepository interface {
	Create(name, password string) error
	Update(id int64, name, password *string) error
	Delete(id int64) error
	Read(id int64) error
}
