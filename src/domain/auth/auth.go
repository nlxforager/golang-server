package auth

type AUTH_MODE string

func (this AUTH_MODE) IsValid() (valid bool) {
	switch this {
	case AUTH_MODE_2FA_PW_E, AUTH_MODE_SIMPLE_PW:
		valid = true
	}
	return
}

const AUTH_MODE_SIMPLE_PW AUTH_MODE = "SIMPLE_PW"
const AUTH_MODE_2FA_PW_E AUTH_MODE = "2FA_PW_E"

type User struct {
	Username string    `json:"username"`
	AuthMode AUTH_MODE `json:"auth_mode"`
}

type ChangeSet struct {
	AuthMode *AUTH_MODE
	Email    *string
}

type AuthService interface {
	GetEmail(username string) (string, error)
	ByPasswordAndUsername(username, password string) (error, *User)
	RegisterUsernamePassword(username, password string) error
	ModifyUser(username string, mode ChangeSet) error

	AuthService_OTP
	JwtService
}

type JwtService interface {
	CreateWeakToken(username string, authMode AUTH_MODE) (string, error)
	CreateStrongToken(username string, authMode AUTH_MODE) (string, error)
	ValidateAndGetClaims(tokenString string) (map[string]string, error)
}

type AuthService_OTP interface {
	VerifyOTP(otp string, s string) error
	SetOTP(username string, otp func() string) error
	OtpGen() string
}
