package service

import (
	"strconv"

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

	return file.GetFilePath(f), nil
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

	return file.GetFilePath(f), nil
}
