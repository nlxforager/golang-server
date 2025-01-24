package auth

import "errors"

type MockUser struct {
	Username string
	Password string
}

type UsernameStore struct {
	UserByUsernames map[string]MockUser
}

func NewStore() *UsernameStore {
	return &UsernameStore{}
}

type MockAuth struct {
	UsernameStore
}

func (m MockAuth) ByPasswordAndUsername(username, password string) error {
	if password == m.UserByUsernames[username].Password {
		return nil
	}
	return errors.New("wrong username or password")
}

func NewMockAuth() *MockAuth {
	return &MockAuth{
		UsernameStore: UsernameStore{
			UserByUsernames: make(map[string]MockUser),
		},
	}
}
