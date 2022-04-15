package file

import (
	"encoding/base64"
	"time"

	"github.com/yzx9/otodo/domain/session"
	"github.com/yzx9/otodo/infrastructure/util"
)

// Configurable
const preSignMaxExpiresIn = 6 * time.Hour

type FilePreSignClaims struct {
	session.JWTClaims

	FileID int64 `json:"fileID"`
}

func ParseFilePreSignID(filePresignedID string) (*File, error) {
	write := func() (*File, error) {
		return nil, util.NewErrorWithPreconditionFailed("invalid file")
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

	file, err := FileRepository.Find(claims.FileID)
	if err != nil {
		return write()
	}

	return file, nil
}

func (file File) CreateFilePreSignID(userID int64, exp int) (string, error) {
	expiresIn := time.Duration(exp * int(time.Second))
	if expiresIn > preSignMaxExpiresIn {
		return "", util.NewErrorWithPreconditionFailed("expires is too long")
	}

	if _, err := GetFileByUser(userID, file.ID); err != nil {
		return "", err
	}

	token := session.NewToken(FilePreSignClaims{
		JWTClaims: session.NewClaims(userID, expiresIn),
		FileID:    file.ID,
	})
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}
