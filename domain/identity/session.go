package identity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// aggregate
type session struct {
	sessionID string
	userID    int64
}

func (s session) UserID() int64 {
	return s.userID
}

// entity
type UserInvalidRefreshToken struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID  int64
	TokenID string
}

func (s session) Deactivate() error {
	model := UserInvalidRefreshToken{
		UserID:  s.userID,
		TokenID: s.sessionID,
	}
	if err := UserInvalidRefreshTokenRepository.Save(&model); err != nil {
		return Error{fmt.Errorf("fails to save user invalid refresh token: %w", err)}
	}

	return nil
}

// value object
type sessionToken struct {
	Token     string
	Type      sessionTokenType
	ExpiresIn int64
}

type sessionTokenType int

const (
	AccessToken sessionTokenType = iota
	RefreshToken
)

// generate access token
func (s session) NewAccessToken() (sessionToken, error) {
	exp := Conf.AccessTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))

	return sessionToken{
		Token:     s.newToken(dur),
		Type:      AccessToken,
		ExpiresIn: int64(exp),
	}, nil
}

// generate refresh token
func (s session) NewRefreshToken(exp int) (sessionToken, error) {
	if exp <= 0 {
		exp = Conf.RefreshTokenExpiresInDefault
	} else if exp > Conf.RefreshTokenExpiresInMax {
		exp = Conf.RefreshTokenExpiresInMax
	}
	dur := time.Duration(exp * int(time.Second))

	return sessionToken{
		Token:     s.newToken(dur),
		Type:      RefreshToken,
		ExpiresIn: int64(exp),
	}, nil
}

func (s session) ShouldRefreshAccessToken(accessToken string) bool {
	claims, err := parseSessionToken(accessToken)
	if err != nil {
		return false
	}

	thd := Conf.AccessTokenRefreshThreshold
	dur := time.Duration(thd * int(time.Second))
	return time.Now().Add(dur).Unix() > claims.ExpiresAt
}

func (s session) newToken(exp time.Duration) string {
	now := time.Now().UTC()
	claims := sessionTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    Conf.TokenIssuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(exp).Unix(),
		},
		UserID:    s.userID,
		SessionID: s.sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Conf.TokenHmacSecret)
	if err != nil {
		return ""
	}

	return tokenString
}

type sessionTokenClaims struct {
	jwt.StandardClaims

	UserID    int64  `json:"uid"`
	SessionID string `json:"sid"`
}

func parseSessionToken(tokenString string) (*sessionTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&sessionTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return Conf.TokenHmacSecret, nil
		})
	if err != nil || !token.Valid {
		return nil, InvalidCredential
	}

	claims := token.Claims.(*sessionTokenClaims)
	return claims, nil
}
