package file

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/util"
)

// Configurable
const preSignMaxExpiresIn = 6 * time.Hour

type FilePreSignClaims struct {
	jwt.StandardClaims

	UserID int64 `json:"uid"`
	FileID int64 `json:"fileID"`
}

func GetFileByPreSignID(filePresignedID string) (*File, error) {
	write := func() (*File, error) {
		return nil, util.NewErrorWithPreconditionFailed("invalid file")
	}

	payload, err := base64.StdEncoding.DecodeString(filePresignedID)
	if err != nil {
		return write()
	}

	token, err := jwt.ParseWithClaims(string(payload), &FilePreSignClaims{}, keyFunc)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if err != nil || !token.Valid {
		return write()
	}

	claims, ok := token.Claims.(*FilePreSignClaims)
	if !ok {
		return write()
	}

	file, err := FileRepository.Find(claims.FileID)
	if err != nil {
		return write()
	}

	return file, nil
}

func (file File) CreateFilePreSignID(userID int64, exp int) (string, error) {
	if !file.CanAccessByUser(userID) {
		return "", fmt.Errorf("permission denied")
	}

	expiresIn := time.Duration(exp * int(time.Second))
	if expiresIn > preSignMaxExpiresIn {
		expiresIn = preSignMaxExpiresIn
	}

	now := time.Now().UTC()
	claims := FilePreSignClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.Secret.TokenIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(expiresIn).Unix(),
		},
		UserID: userID,
		FileID: file.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Secret.TokenHmacSecret)
	if err != nil {
		return "", nil
	}

	tokenEncoded := base64.StdEncoding.EncodeToString([]byte(tokenString))
	return tokenEncoded, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return config.Secret.TokenHmacSecret, nil
}
