package mocks

import (
	"errors"
	"time"

	"golang-server/src/domain/authservice"
	"golang-server/src/domain/authservice/jwt"
	"golang-server/src/domain/authservice/otp"
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
	Generator otp.SimpleGenerator
}

func (m MockAuth) ValidateAndGetClaims(tokenString string) (map[string]string, error) {
	service := mockJwtService()
	claims, err := service.GetClaims(tokenString)

	if err != nil {
		return nil, err
	}

	cmap := make(map[string]string)

	for k, v := range claims {
		vv, _ := v.(string)
		cmap[k] = vv
	}
	return cmap, nil
}

func (m MockAuth) ModifyUser(username string, set authservice.ChangeSet) error {
	u, ok := m.UserByUsernames[username]
	if !ok {
		return errors.New("user not found")
	}
	if set.AuthMode != nil {
		u.Mode = string(*set.AuthMode)
	}
	if set.Email != nil {
		u.Email = *set.Email
	}
	m.UserByUsernames[username] = u
	return nil
}

func (m MockAuth) RegisterUsernamePassword(username, password string) error {
	m.UserByUsernames[username] = MockUser{
		Username: username,
		Password: password,
		Email:    "",
		Mode:     string(authservice.AUTH_MODE_SIMPLE_PW),
	}

	return nil
}

func mockJwtService() jwt.Service {
	return jwt.Service{
		SecretKey: "TEST",
		Issuer:    "TEST-MOCK-SERVER",
	}
}

func (m MockAuth) CreateWeakToken(username string, authMode authservice.AUTH_MODE) (string, error) {
	service := mockJwtService()
	return service.CreateWeakToken(username, map[string]string{
		"auth_mode": string(authMode),
	})
}

func (m MockAuth) CreateStrongToken(username string, authMode authservice.AUTH_MODE) (string, error) {
	service := mockJwtService()
	return service.CreateWeakToken(username, map[string]string{
		"auth_mode": string(authMode),
		"is_auth":   "true",
	})
}

func (m MockAuth) GetClaimFromToken(username string, claimKey string) (string, error) {
	service := mockJwtService()
	return service.GetClaimFromToken(username, claimKey)
}

func (m MockAuth) GetUsernameFromToken(username string) (string, error) {
	service := jwt.Service{
		SecretKey: "TEST",
		Issuer:    "TEST-MOCK-SERVER",
	}

	return service.GetUsernameFromToken(username)
}

func (m MockAuth) GetEmail(username string) (string, error) {
	email := m.UserByUsernames[username].Email
	if email == "" {
		return "", errors.New("email not found")
	}
	return email, nil
}

func (m MockAuth) OtpGen() func() string {
	return func() string { return "999999" }
}

func (m MockAuth) SetOTP(username string, otp otp.OtpGen) (string, error) {
	o := otp()
	m.OtpByUsername[username] = Otp{Value: o, Expiry: time.Now().Add(time.Minute)}
	return o, nil
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

func (m MockAuth) ByPasswordAndUsername(username, password string) (error, *authservice.User) {
	u := m.UserByUsernames[username]
	if password == u.Password {
		return nil, &authservice.User{
			Username: username,
			AuthMode: authservice.AUTH_MODE(u.Mode),
		}
	}

	return errors.New("wrong username or password"), nil
}

var _ authservice.AuthService = (*MockAuth)(nil)

func NewMockAuth() *MockAuth {
	return &MockAuth{
		UsernameStore: UsernameStore{
			UserByUsernames: make(map[string]MockUser),
		},
		OtpStore: OtpStore{
			OtpByUsername: make(map[string]Otp),
		},

		Generator: otp.SimpleGenerator{},
	}
}
