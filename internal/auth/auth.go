package auth

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %s", err)
	}
	return string(hashedPass), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetBearerToken(headers http.Header) (string, error) {
	authValue := headers.Get("Authorization")
	if authValue == "" {
		return "", fmt.Errorf("no authorization header found")
	}

	strippedAuth := strings.TrimPrefix(authValue, "Bearer ")

	return strippedAuth, nil
}
