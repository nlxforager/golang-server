package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	SecretKey string
	Issuer    string
}

// Function to create Service tokens with claims
func (a *Service) CreateWeakToken(username string, claims map[string]string) (string, error) {
	// Create a new Service token with claims

	mclaim := jwt.MapClaims{
		"sub": username, // Subject (user identifier)
		"iss": a.Issuer, // Issuer
		//"aud": getRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	}

	for k, v := range claims {
		mclaim[k] = v
	}
	return a.signToken(mclaim)
}

// Function to create Service tokens with claims
func (a *Service) GetUsernameFromToken(tokenString string) (string, error) {
	return a.GetClaimFromToken(tokenString, "sub")
}

// Function to create Service tokens with claims
func (a *Service) GetAuthModeFromToken(tokenString string) (string, error) {
	return a.GetClaimFromToken(tokenString, "auth_mode")
}

// Function to create Service tokens with claims
func (a *Service) GetClaimFromToken(tokenString string, claim string) (string, error) {
	maps, err := a.GetClaims(tokenString)
	if err != nil {
		return "", err
	}
	field, _ := maps[claim].(string)
	return field, nil
}

func (a *Service) signToken(claims jwt.MapClaims) (string, error) {
	// Print information about the created token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Service) GetClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	maps, _ := token.Claims.(jwt.MapClaims)
	return maps, nil
}
