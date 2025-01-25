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
func (a *Service) CreateTokenUsernameOnly(username string) (string, error) {
	// Create a new Service token with claims
	return a.signToken(jwt.MapClaims{
		"sub": username, // Subject (user identifier)
		"iss": a.Issuer, // Issuer
		//"aud": getRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})
}

// Function to create Service tokens with claims
func (a *Service) GetUsernameFromToken(tokenString string) (string, error) {
	// Create a new Service token with claims
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.SecretKey), nil
	})

	// Check for verification errors
	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	maps, _ := token.Claims.(jwt.MapClaims)
	username, _ := maps["sub"].(string)
	// Return the verified token
	return username, nil
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
