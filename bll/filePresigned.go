package bll

import (
	"encoding/base64"
	"time"

	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/utils"
)

// Configurable
var fileSignedExpress = 6 * time.Hour

type filePresignedPayload struct {
	TokenClaims
	UserID string `json:"uid"`
	FileID string `json:"todo_id"`
}

func CreateFilePresignedID(userID, fileID string) string {
	payload := filePresignedPayload{
		TokenClaims: NewClaims(userID, fileSignedExpress),
		UserID:      userID,
		FileID:      fileID,
	}
	token := NewToken(payload)
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func ParseFileSignedID(filePresignedID string) (string, error) {
	write := func() (string, error) {
		return "", utils.NewError(otodo.ErrAbort, "invalid presigned file id")
	}

	payload, err := base64.StdEncoding.DecodeString(filePresignedID)
	if err != nil {
		return write()
	}

	token, err := ParseToken(string(payload), &filePresignedPayload{})
	if err != nil || !token.Valid {
		return write()
	}

	claims, ok := token.Claims.(filePresignedPayload)
	if !ok {
		return write()
	}

	return claims.FileID, nil
}
