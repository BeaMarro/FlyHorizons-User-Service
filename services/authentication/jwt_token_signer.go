package authentication

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type JwtTokenSigner struct{}

func NewJwtTokenSigner() *JwtTokenSigner {
	return &JwtTokenSigner{}
}

func (s *JwtTokenSigner) SignToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	fmt.Println("Secret Key", secretKey)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
