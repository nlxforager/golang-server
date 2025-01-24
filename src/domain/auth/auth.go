package auth

type AuthService interface {
	ByPasswordAndUsername(username, password string) error
}
