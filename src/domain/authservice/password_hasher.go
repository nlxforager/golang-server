package authservice

type PasswordHasher struct {
	secret string
}

func NewPasswordHasher(secret string) *PasswordHasher {
	return &PasswordHasher{
		secret: secret,
	}
}
func (p PasswordHasher) HashedPassword(password string) string {
	return p.secret + password
}
