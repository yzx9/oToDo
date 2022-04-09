package service

import (
	"strconv"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/infrastructure/util"
)

func GetFilePath(userID int64, fileID string) (string, error) {
	id, err := strconv.ParseInt(fileID, 10, 64)
	if err != nil {
		return GetPreSignFilePath(fileID)
	}

	f, err := file.OwnFile(userID, id)
	if err != nil {
		return "", err
	}

	return f.GetFilePath(), nil
}

func GetPreSignFilePath(fileID string) (string, error) {
	id, err := file.ParseFilePreSignID(fileID)
	if err != nil {
		return "", util.NewErrorWithNotFound("file not found")
	}

	f, err := file.GetFile(id)
	if err != nil {
		return "", err
	}

	return f.GetFilePath(), nil
}

func PreSignFile(payload dto.FilePreSign) (dto.FilePreSignResult, error) {
	presigned, err := file.CreateFilePreSignID(payload.UserID, payload.FileID, payload.ExpiresIn)
	if err != nil {
		return dto.FilePreSignResult{}, err
	}

	return dto.FilePreSignResult{
		FileID: presigned,
	}, nil
}
