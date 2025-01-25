package otp

type OtpGen = func() string
type OtpMockGenerator struct{}

func (*OtpMockGenerator) Generate() string {
	return "123456"
}
