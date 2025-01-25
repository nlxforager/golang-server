package auth

import (
	"errors"
	"time"

	"golang-server/src/domain/auth/jwt"
	"golang-server/src/domain/auth/otp"
)

type MockUser struct {
	Username string
	Password string
	Email    string
	Mode     string
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

func (m MockAuth) RegisterUsernamePassword(username, password string) error {
	m.UserByUsernames[username] = MockUser{
		Username: username,
		Password: password,
		Email:    "",
		Mode:     string(AUTH_MODE_SIMPLE_PW),
	}

	return nil
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

	username, err := service.GetUsernameFromToken(token)
	if err != nil {
		return err
	}

	otp := m.OtpByUsername[username]
	if otpVal != otp.Value {
		return errors.New("invalid otp")
	}
	return nil
}

func (m MockAuth) ByPasswordAndUsername(username, password string) (error, *User) {
	u := m.UserByUsernames[username]
	if password == u.Password {
		return nil, &User{
			Username: username,
			AuthMode: AUTH_MODE(u.Mode),
		}
	}

	return errors.New("wrong username or password"), nil
}

var _ AuthService = (*MockAuth)(nil)

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
