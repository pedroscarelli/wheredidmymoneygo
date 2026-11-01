package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenDuration = 2 * time.Hour

var ErrMissingJwtSecret = errors.New("JWT_SECRET_KEY is not set")

func JwtSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return nil, ErrMissingJwtSecret
	}

	return []byte(secret), nil
}

func GenerateJwtToken(accountID string) (string, error) {
	secret, err := JwtSecret()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": accountID,
		"exp": time.Now().Add(tokenDuration).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}
