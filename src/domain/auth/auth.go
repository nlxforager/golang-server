package auth

type AUTH_MODE string

const AUTH_MODE_SIMPLE_PW AUTH_MODE = "SIMPLE_PW"
const AUTH_MODE_2FA_PW_E AUTH_MODE = "2FA_PW_E"

type User struct {
	Username string    `json:"username"`
	AuthMode AUTH_MODE `json:"auth_mode"`
}

type AuthService interface {
	GetEmail(username string) (string, error)
	ByPasswordAndUsername(username, password string) (error, *User)

	AuthService_OTP
	JwtService
}

type JwtService interface {
	CreateTokenUsernameOnly(username string) (string, error)
}

type AuthService_OTP interface {
	VerifyOTP(otp string, s string) error
	SetOTP(username string, otp func() string) error
	OtpGen() string
}
