package interfaces

import "github.com/golang-jwt/jwt"

type TokenSigner interface {
	SignToken(claims jwt.Claims) (string, error)
}
