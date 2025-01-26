package auth

import (
	"context"
	"errors"
	"time"

	"golang-server/src/domain/auth/jwt"
	"golang-server/src/domain/auth/otp"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	Repo         Repository
	Redis        *redis.Client
	JwtIssuer    string
	JwtSecretKey string
}

var _ AuthService = (*Service)(nil)

func NewService(repository *Repository, redisa *redis.Client, issuer string, secret string) (*Service, error) {
	if repository == nil || redisa == nil {
		return nil, errors.New("persistence is nil")
	}

	if issuer == "" || secret == "" {
		return nil, errors.New("jwt config fail.")
	}

	return &Service{
		Repo:         Repository{},
		Redis:        redisa,
		JwtIssuer:    issuer,
		JwtSecretKey: secret,
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

func (s Service) SetOTP(username string, otp func() string) (string, error) {
	o := otp()
	rdb := s.Redis
	st := rdb.Set(context.TODO(), "set_otp/"+username, o, time.Second*10)
	err := st.Err()
	if err != nil {
		return "", err
	}
	return o, nil
}

// VerifyOTP
// do not check if is weak or strong token, that is OTP may be the first factor of a MFA flow.
func (s Service) VerifyOTP(otp string, token string) error {
	//TODO implement me
	tokenS, err := s.ValidateAndGetClaims(token)
	if err != nil {
		return err
	}

	username := tokenS["sub"]
	val := s.Redis.Get(context.TODO(), "set_otp/"+username)

	err = val.Err()
	if err != nil {
		return err
	}

	vv, _ := val.Result()
	if vv != otp {
		return errors.New("otp verification failed")
	}
	return nil
}

func (s Service) OtpGen() func() string {
	generator := otp.SimpleGenerator{}
	return func() string { return generator.Generate() }
}

func (s Service) CreateWeakToken(username string, authMode AUTH_MODE) (string, error) {
	service := jwt.Service{
		SecretKey: s.JwtSecretKey,
		Issuer:    s.JwtIssuer,
	}

	return service.CreateWeakToken(username, map[string]string{
		"auth_mode": string(authMode),
	})
}

func (s Service) CreateStrongToken(username string, authMode AUTH_MODE) (string, error) {
	service := jwt.Service{
		SecretKey: s.JwtSecretKey,
		Issuer:    s.JwtIssuer,
	}

	return service.CreateWeakToken(username, map[string]string{
		"auth_mode": string(authMode),
		"is_auth":   "true",
	})
}

func (s Service) ValidateAndGetClaims(tokenString string) (map[string]string, error) {
	service := jwt.Service{
		SecretKey: s.JwtSecretKey,
		Issuer:    s.JwtIssuer,
	}
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
