package file

import (
	"encoding/base64"
	"time"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/user"
	"github.com/yzx9/otodo/infrastructure/util"
)

// Configurable
const fileSignedMaxExpiresIn = 6 * time.Hour

func CreateFilePreSignID(userID, fileID int64) (string, error) {
	const max = int(fileSignedMaxExpiresIn / time.Second)
	return CreateFilePreSignIDWithExp(userID, fileID, max)
}

func CreateFilePreSignIDWithExp(userID, fileID int64, exp int) (string, error) {
	expiresIn := time.Duration(exp * int(time.Second))
	if expiresIn > fileSignedMaxExpiresIn {
		return "", util.NewErrorWithPreconditionFailed("expires is too long")
	}

	_, err := OwnFile(userID, fileID)
	if err != nil {
		return "", err
	}

	token := user.NewToken(dto.FilePreSignClaims{
		TokenClaims: user.NewClaims(userID, expiresIn),
		UserID:      userID,
		FileID:      fileID,
	})
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

func GetPreSignFilePath(fileID string) (string, error) {
	id, err := parseFilePreSignID(fileID)
	if err != nil {
		return "", util.NewErrorWithNotFound("file not found")
	}

	file, err := GetFile(id)
	if err != nil {
		return "", err
	}

	return getFilePath(file), nil
}

func parseFilePreSignID(filePresignedID string) (int64, error) {
	write := func() (int64, error) {
		return 0, util.NewErrorWithPreconditionFailed("invalid presigned file id")
	}

	payload, err := base64.StdEncoding.DecodeString(filePresignedID)
	if err != nil {
		return write()
	}

	token, err := user.ParseToken(string(payload), &dto.FilePreSignClaims{})
	if err != nil || !token.Valid {
		return write()
	}

	claims, ok := token.Claims.(*dto.FilePreSignClaims)
	if !ok {
		return write()
	}

	return claims.FileID, nil
}
