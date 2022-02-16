package dto

import (
	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	jwt.StandardClaims

	UserID string `json:"uid"`
}
