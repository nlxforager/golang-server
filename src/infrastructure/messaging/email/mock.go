package email

import (
	"context"
	"fmt"
	"time"
)

type MockOtpSingleSendReceiver struct {
	c chan string
}

func (m *MockOtpSingleSendReceiver) OTP(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case msg := <-m.c:
		return msg, nil
	}
}

func (m *MockOtpSingleSendReceiver) SendOTP(email, otp string) error {

	fmt.Printf("sending... otp")
	<-time.After(2 * time.Second)
	fmt.Printf("sent... otp")
	m.c <- otp
	return nil
}

func NewMockOtpSingleSendReceiver() *MockOtpSingleSendReceiver {
	return &MockOtpSingleSendReceiver{
		c: make(chan string),
	}
}
