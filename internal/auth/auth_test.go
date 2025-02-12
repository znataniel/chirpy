package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	duration := time.Minute * 10

	tokStr, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("failure creating jwt: %v", err)
	}

	validatedUID, err := ValidateJWT(tokStr, secret)
	if err != nil {
		t.Fatalf("failure validating jwt: %v", err)
	}

	if validatedUID != userID {
		t.Fatalf("uuid do not match:\n- %s\n- %s", userID, validatedUID)
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	altSecret := "anothersecret"
	duration := time.Minute * 10

	tokStr, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("failure creating jwt: %v", err)
	}

	_, err = ValidateJWT(tokStr, altSecret)
	if err == nil {
		t.Fatalf("jwt should not have validated with secret: %v", altSecret)
	}

}

func TestValidateExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "testsecret"
	duration := time.Millisecond * 100

	tokStr, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("failure creating jwt: %v", err)
	}

	_, err = ValidateJWT(tokStr, secret)
	if err == nil {
		t.Fatalf("jwt should have expired after: %s", duration)
	}

}
