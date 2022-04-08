package file

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
)

// Configurable
const fileSignedMaxExpiresIn = 6 * time.Hour

func CreateFilePreSignID(userID, fileID int64, exp int) (string, error) {
	expiresIn := time.Duration(exp * int(time.Second))
	if expiresIn > fileSignedMaxExpiresIn {
		return "", util.NewErrorWithPreconditionFailed("expires is too long")
	}

	_, err := OwnFile(userID, fileID)
	if err != nil {
		return "", err
	}

	// TODO[bug]: following code make cycle dep
	// token := user.NewToken(dto.FilePreSignClaims{
	// 	TokenClaims: user.NewClaims(userID, expiresIn),
	// 	UserID:      userID,
	// 	FileID:      fileID,
	// })
	token := ""
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

func ParseFilePreSignID(filePresignedID string) (int64, error) {
	return 0, fmt.Errorf("TODO: resolve cycle dep")

	// write := func() (int64, error) {
	// 	return 0, util.NewErrorWithPreconditionFailed("invalid presigned file id")
	// }

	// payload, err := base64.StdEncoding.DecodeString(filePresignedID)
	// if err != nil {
	// 	return write()
	// }

	// TODO[bug]: following code make cycle dep
	// token, err := user.ParseToken(string(payload), &dto.FilePreSignClaims{})
	// if err != nil || !token.Valid {
	// 	return write()
	// }

	// claims, ok := token.Claims.(*dto.FilePreSignClaims)
	// if !ok {
	// 	return write()
	// }

	// return claims.FileID, nil
}
