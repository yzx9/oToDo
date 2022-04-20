package identity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/config"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/util"
)

type Session struct {
	SessionID string
	UserID    int64
}

func LoginByCredential(userName, password string) (Session, error) {
	user, err := UserRepository.FindByUserName(userName)
	if err != nil ||
		user.Password == nil ||
		!user.IsSamePassword(password) {
		return Session{}, util.NewErrorWithBadRequest("invalid credential")
	}

	return Session{
		SessionID: uuid.NewString(),
		UserID:    user.ID,
	}, nil
}

func LoginByGithubOAuth(code, state string) (Session, error) {
	oauth, err := GetOAuthEntryByState(state)
	if err != nil {
		return Session{}, fmt.Errorf("invalid state")
	}

	user, err := oauth.GetUserByGithub(code)
	if err != nil {
		return Session{}, fmt.Errorf("fails to get user: %w", err)
	}

	return Session{
		SessionID: uuid.NewString(),
		UserID:    user.ID,
	}, nil
}

func LoginByAccessToken(token string) (Session, error) {
	claims, err := parseSessionToken(token)
	if err != nil {
		return Session{}, util.NewError(errors.ErrUnauthorized, "invalid token")
	}

	return Session{
		SessionID: claims.RefreshTokenID,
		UserID:    claims.UserID,
	}, nil
}

func LoginByRefreshToken(token string) (Session, error) {
	claims, err := parseSessionToken(token)
	if err != nil {
		return Session{}, util.NewError(errors.ErrUnauthorized, "invalid token")
	}

	valid, err := UserInvalidRefreshTokenRepository.Exist(claims.UserID, claims.Id)
	if err != nil || valid {
		return Session{}, fmt.Errorf("invalid refresh token")
	}

	return Session{
		SessionID: claims.Id,
		UserID:    claims.UserID,
	}, nil
}

func (s Session) Logout() error {
	_, err := CreateUserInvalidRefreshToken(s.UserID, s.SessionID)
	return err
}

type Token struct {
	Token     string
	Type      TokenType
	ExpiresIn int64
}

// generate access token
func (s Session) NewAccessToken() (Token, error) {
	user, err := UserRepository.Find(s.UserID)
	if err != nil {
		return Token{}, fmt.Errorf("fails to get user, %w", err)
	}

	exp := config.Session.AccessTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))

	return Token{
		Token:     newClaims(user.ID, s.SessionID, dur).toToken(),
		Type:      AccessToken,
		ExpiresIn: int64(exp),
	}, nil
}

// generate refresh token
func (s Session) NewRefreshToken(exp int) (Token, error) {
	if exp <= 0 {
		exp = config.Session.RefreshTokenExpiresInDefault
	} else if exp > config.Session.RefreshTokenExpiresInMax {
		exp = config.Session.RefreshTokenExpiresInMax
	}
	dur := time.Duration(exp * int(time.Second))

	claims := newClaims(s.UserID, uuid.NewString(), dur)
	claims.Id = claims.RefreshTokenID
	return Token{
		Token:     claims.toToken(),
		Type:      RefreshToken,
		ExpiresIn: int64(exp),
	}, nil
}

func (s Session) ShouldRefreshAccessToken(accessToken string) bool {
	claims, err := parseSessionToken(accessToken)
	if err != nil {
		return false
	}

	thd := config.Session.AccessTokenRefreshThreshold
	dur := time.Duration(thd * int(time.Second))
	return time.Now().Add(dur).Unix() > claims.ExpiresAt
}

type sessionTokenClaims struct {
	jwt.StandardClaims

	UserID         int64  `json:"uid"`
	RefreshTokenID string `json:"rti,omitempty"`
}

func parseSessionToken(tokenString string) (*sessionTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&sessionTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.Secret.TokenHmacSecret, nil
		})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims := token.Claims.(*sessionTokenClaims)
	return claims, nil
}

func newClaims(userID int64, sessionID string, exp time.Duration) sessionTokenClaims {
	now := time.Now().UTC()
	return sessionTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.Secret.TokenIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(exp).Unix(),
		},
		UserID:         userID,
		RefreshTokenID: sessionID,
	}
}

func (claims sessionTokenClaims) toToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Secret.TokenHmacSecret)
	if err != nil {
		return ""
	}

	return tokenString
}
