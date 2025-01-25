package auth

import (
	"errors"
	"fmt"
	"time"

	"golang-server/src/domain/auth/jwt"
	"golang-server/src/domain/auth/otp"
)

type MockUser struct {
	Username string
	Password string
	Email    string
}

type UsernameStore struct {
	UserByUsernames map[string]MockUser
}

func NewStore() *UsernameStore {
	return &UsernameStore{}
}

type Otp struct {
	Value  string
	Expiry time.Time
}

type OtpStore struct {
	OtpByUsername map[string]Otp
}

type MockAuth struct {
	UsernameStore
	OtpStore
	Generator otp.OtpMockGenerator
}

func (m MockAuth) CreateTokenUsernameOnly(username string) (string, error) {
	service := jwt.Service{
		SecretKey: "TEST",
		Issuer:    "TEST-MOCK-SERVER",
	}
	return service.CreateTokenUsernameOnly(username)
}

func (m MockAuth) GetUsernameFromToken(username string) (string, error) {
	service := jwt.Service{
		SecretKey: "TEST",
		Issuer:    "TEST-MOCK-SERVER",
	}

	return service.CreateTokenUsernameOnly(username)
}

func (m MockAuth) GetEmail(username string) (string, error) {
	email := m.UserByUsernames[username].Email
	if email == "" {
		return "", errors.New("email not found")
	}
	return email, nil
}

func (m MockAuth) OtpGen() string {
	return "999999"
}

func (m MockAuth) SetOTP(username string, otp otp.OtpGen) error {
	m.OtpByUsername[username] = Otp{Value: otp(), Expiry: time.Now().Add(time.Minute)}
	return nil
}

func (m MockAuth) VerifyOTP(otpVal string, token string) error {
	service := jwt.Service{
		SecretKey: "TEST",
		Issuer:    "TEST-MOCK-SERVER",
	}
	fmt.Printf("token.username is %s\n", token)

	username, err := service.GetUsernameFromToken(token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	otp := m.OtpByUsername[username]
	if otpVal != otp.Value {
		return errors.New("invalid otp")
	}
	return nil
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
		OtpStore: OtpStore{
			OtpByUsername: make(map[string]Otp),
		},

		Generator: otp.OtpMockGenerator{},
	}
}
