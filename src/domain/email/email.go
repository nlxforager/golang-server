package email

import (
	"log"
	"net/smtp"
)

type OTPEmailer interface {
	SendOTP(email string, otp string) error
}

type ClientSimpleService struct {
	Password string
	Email    string
}

func NewSimpleClientService(email string, password string) (*ClientSimpleService, error) {
	return &ClientSimpleService{
		Email:    email,
		Password: password,
	}, nil
}

func (service *ClientSimpleService) SendOTP(to string, otp string) error {
	// Create a new message
	from := service.Email
	pass := service.Password

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		otp

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}
	return nil
}
