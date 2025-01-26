package otp

import "math/rand"

type OtpGen = func() string
type SimpleGenerator struct{}

var letterRunes = []rune("123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (*SimpleGenerator) Generate() string {
	return RandStringRunes(6)
}
