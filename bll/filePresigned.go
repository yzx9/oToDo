package bll

import (
	"encoding/base64"
	"time"

	"github.com/yzx9/otodo/utils"
)

// Configurable
const fileSignedMaxExpiresIn = 6 * time.Hour

type filePresignedPayload struct {
	TokenClaims
	UserID string `json:"uid"`
	FileID string `json:"fileID"`
}

func CreateFilePresignedID(userID, fileID string) (string, error) {
	const max = int(fileSignedMaxExpiresIn / time.Second)
	return CreateFilePresignedIDWithExp(userID, fileID, max)
}

func CreateFilePresignedIDWithExp(userID, fileID string, exp int) (string, error) {
	expiresIn := time.Duration(exp * int(time.Second))
	if expiresIn > fileSignedMaxExpiresIn {
		return "", utils.NewErrorWithPreconditionFailed("expires is too long")
	}

	_, err := OwnFile(userID, fileID)
	if err != nil {
		return "", err
	}

	token := NewToken(filePresignedPayload{
		TokenClaims: NewClaims(userID, expiresIn),
		UserID:      userID,
		FileID:      fileID,
	})
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

func ParseFileSignedID(filePresignedID string) (string, error) {
	write := func() (string, error) {
		return "", utils.NewErrorWithPreconditionFailed("invalid presigned file id")
	}

	payload, err := base64.StdEncoding.DecodeString(filePresignedID)
	if err != nil {
		return write()
	}

	token, err := ParseToken(string(payload), &filePresignedPayload{})
	if err != nil || !token.Valid {
		return write()
	}

	claims, ok := token.Claims.(*filePresignedPayload)
	if !ok {
		return write()
	}

	return claims.FileID, nil
}
