package email

type EmailService interface {
	SendOTP(email string, otp string) error
}
