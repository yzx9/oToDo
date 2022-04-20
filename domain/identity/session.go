package identity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// aggregate
type Session struct {
	sessionID string
	userID    int64
}

func (s Session) UserID() int64 {
	return s.userID
}

func LoginByCredential(userName, password string) (Session, error) {
	user, err := UserRepository.FindByUserName(userName)
	if err != nil ||
		user.Password == nil ||
		!user.ValidatePassword(password) {
		return Session{}, InvalidCredential
	}

	return Session{
		userID:    user.ID,
		sessionID: uuid.NewString(),
	}, nil
}

func LoginByGithubOAuth(code, state string) (Session, error) {
	oauth, err := GetOAuthEntryByState(state)
	if err != nil {
		return Session{}, InvalidCredential
	}

	user, err := oauth.GetUserByGithub(code)
	if err != nil {
		return Session{}, newErr(fmt.Errorf("fails to get user: %w", err))
	}

	return Session{
		userID:    user.ID,
		sessionID: uuid.NewString(),
	}, nil
}

func LoginByAccessToken(token string) (Session, error) {
	claims, err := parseSessionToken(token)
	if err != nil {
		return Session{}, InvalidCredential
	}

	return Session{
		userID:    claims.UserID,
		sessionID: claims.SessionID,
	}, nil
}

func LoginByRefreshToken(token string) (Session, error) {
	claims, err := parseSessionToken(token)
	if err != nil {
		return Session{}, InvalidCredential
	}

	valid, err := UserInvalidRefreshTokenRepository.Exist(claims.UserID, claims.Id)
	if err != nil || valid {
		return Session{}, InvalidCredential
	}

	return Session{
		userID:    claims.UserID,
		sessionID: claims.Id,
	}, nil
}

// entity
type UserInvalidRefreshToken struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID  int64
	TokenID string
}

func (s Session) Inactive() error {
	model := UserInvalidRefreshToken{
		UserID:  s.userID,
		TokenID: s.sessionID,
	}
	if err := UserInvalidRefreshTokenRepository.Save(&model); err != nil {
		return newErr(fmt.Errorf("fails to save user invalid refresh token: %w", err))
	}

	return nil
}

// value object
type Token struct {
	Token     string
	Type      TokenType
	ExpiresIn int64
}

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

// generate access token
func (s Session) NewAccessToken() (Token, error) {
	exp := Conf.AccessTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))

	return Token{
		Token:     s.newToken(dur),
		Type:      AccessToken,
		ExpiresIn: int64(exp),
	}, nil
}

// generate refresh token
func (s Session) NewRefreshToken(exp int) (Token, error) {
	if exp <= 0 {
		exp = Conf.RefreshTokenExpiresInDefault
	} else if exp > Conf.RefreshTokenExpiresInMax {
		exp = Conf.RefreshTokenExpiresInMax
	}
	dur := time.Duration(exp * int(time.Second))

	return Token{
		Token:     s.newToken(dur),
		Type:      RefreshToken,
		ExpiresIn: int64(exp),
	}, nil
}

func (s Session) ShouldRefreshAccessToken(accessToken string) bool {
	claims, err := parseSessionToken(accessToken)
	if err != nil {
		return false
	}

	thd := Conf.AccessTokenRefreshThreshold
	dur := time.Duration(thd * int(time.Second))
	return time.Now().Add(dur).Unix() > claims.ExpiresAt
}

func (s Session) newToken(exp time.Duration) string {
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
