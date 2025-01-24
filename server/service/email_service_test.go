package service

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerateVerificationToken(t *testing.T) {
	email := "test@example.com"
	token, err := generateVerificationToken(email)

	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatal("Token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to parse claims")
	}

	if claims["email"] != email {
		t.Errorf("Expected email %s, got %s", email, claims["email"])
	}

	expirationTime := claims["exp"].(float64)
	if time.Now().Unix() > int64(expirationTime) {
		t.Fatal("Token has expired")
	}
}
