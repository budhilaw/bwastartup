package domain

import "github.com/golang-jwt/jwt/v4"

type JWTService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
