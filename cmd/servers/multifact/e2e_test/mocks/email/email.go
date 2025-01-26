package email

import (
	"context"
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
	<-time.After(2 * time.Second)
	m.c <- otp
	return nil
}

func NewMockOtpSingleSendReceiver() *MockOtpSingleSendReceiver {
	return &MockOtpSingleSendReceiver{
		c: make(chan string),
	}
}
