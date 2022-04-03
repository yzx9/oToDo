package dto

import (
	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	jwt.StandardClaims

	UserID int64 `json:"uid"`
}
