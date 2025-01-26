package authservice

import (
	"context"
	"errors"
	"time"
)
import "github.com/redis/go-redis/v9"

type Service struct {
	Repo  Repository
	Redis *redis.Client
}

var _ AuthService = (*Service)(nil)

func NewService(repository *Repository, redisa *redis.Client) (*Service, error) {
	if repository == nil || redisa == nil {
		return nil, errors.New("persistence is nil")
	}
	return &Service{
		Repo:  Repository{},
		Redis: redisa,
	}, nil
}

const DEFAULT_AUTH AUTH_MODE = AUTH_MODE_SIMPLE_PW

func (s Service) RegisterUsernamePassword(username, plain string) error {
	hashed := PasswordHasher{}.HashedPassword(plain)
	_, err := s.Repo.CreateUser(context.TODO(), CreateUserParams{
		Username:       username,
		HashedPassword: hashed,
		AuthMode:       DEFAULT_AUTH,
	})
	return err
}

func (s Service) ByPasswordAndUsername(username, plain string) (error, *User) {
	hashed := PasswordHasher{}.HashedPassword(plain)
	user, err := s.Repo.GetUser(context.TODO(), GetUserParams{
		Username:       username,
		HashedPassword: hashed,
	})
	if err != nil {
		return err, nil
	}

	return nil, &User{
		Username: user.Username,
		AuthMode: AUTH_MODE(user.AuthMode),
	}
}

func (s Service) GetEmail(username string) (string, error) {
	emails, err := s.Repo.GetEmail(context.TODO(), username)

	if err != nil {
		return "", err
	}

	if len(emails) == 0 {
		return "", errors.New("email not found")
	}

	if len(emails) > 1 {
		return "", errors.New("multiple emails found, service does not support as of now")
	}

	v, _ := emails[0].Value()
	if v == nil {
		return "", errors.New("email not found")
	}

	return v.(string), nil
}

func (s Service) ModifyUser(username string, mode ChangeSet) error {
	var ss *string
	if mode.AuthMode != nil {
		ss = (*string)(mode.AuthMode)
	}
	return s.Repo.UpdateUser(context.TODO(), UpdateUserAuthModeParams{
		AuthMode: ss,
		Email:    mode.Email,
		Username: username,
	})
}

type OTP string

func (o OTP) Otp() string {
	return string(o)
}

func (s Service) SetOTP(username string, otp func() string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})

	s.Redis = rdb

	st := rdb.Set(context.TODO(), "set_otp/"+username, otp(), time.Second*10)
	return st.Err()
}

func (s Service) VerifyOTP(otp string, str string) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) OtpGen() string {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreateWeakToken(username string, authMode AUTH_MODE) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreateStrongToken(username string, authMode AUTH_MODE) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) ValidateAndGetClaims(tokenString string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}
