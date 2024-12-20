package memdb

type Store struct {
}

func New() (*Store, error) {
	return &Store{}, nil
}
