package file

import (
	"encoding/base64"
	"time"

	"github.com/yzx9/otodo/domain/session"
	"github.com/yzx9/otodo/infrastructure/util"
)

// Configurable
const fileSignedMaxExpiresIn = 6 * time.Hour

type FilePreSignClaims struct {
	session.JWTClaims

	FileID int64 `json:"fileID"`
}

func CreateFilePreSignID(userID, fileID int64, exp int) (string, error) {
	expiresIn := time.Duration(exp * int(time.Second))
	if expiresIn > fileSignedMaxExpiresIn {
		return "", util.NewErrorWithPreconditionFailed("expires is too long")
	}

	if _, err := OwnFile(userID, fileID); err != nil {
		return "", err
	}

	token := session.NewToken(FilePreSignClaims{
		JWTClaims: session.NewClaims(userID, expiresIn),
		FileID:    fileID,
	})
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

func ParseFilePreSignID(filePresignedID string) (int64, error) {
	write := func() (int64, error) {
		return 0, util.NewErrorWithPreconditionFailed("invalid presigned file id")
	}

	payload, err := base64.StdEncoding.DecodeString(filePresignedID)
	if err != nil {
		return write()
	}

	token, err := session.ParseToken(string(payload), &FilePreSignClaims{})
	if err != nil || !token.Valid {
		return write()
	}

	claims, ok := token.Claims.(*FilePreSignClaims)
	if !ok {
		return write()
	}

	return claims.FileID, nil
}
