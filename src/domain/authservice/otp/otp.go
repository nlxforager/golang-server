package otp

import "math/rand"

type OtpGen = CreateWeakToken
type OtpMockGenerator struct{}

var letterRunes = []rune("123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func (*OtpMockGenerator) Generate() string {
	return RandStringRunes(6)
}
