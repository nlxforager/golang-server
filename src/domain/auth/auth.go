package auth

type AuthService interface {
	GetEmail(username string) (string, error)
	ByPasswordAndUsername(username, password string) error

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
