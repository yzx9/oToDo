package bll

import (
	"errors"
	"mime/multipart"
	"strings"
	"time"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/utils"
)

const maxFileSize = 8 << 20 // 8MiB

func UploadFile(file *multipart.FileHeader) (string, error) {
	// router.MaxMultipartMemory = MaxFileSize // 限制 Gin 上传文件时最大内存 (默认 32 MiB)
	if file.Size > maxFileSize {
		return "", errors.New("file too large")
	}

	destTemplate, err := dal.GetFileDestTemplate()
	if err != nil {
		return "", errors.New("internal error")
	}

	dest := strings.ReplaceAll(destTemplate, ":filename", file.Filename)
	err = utils.SaveFile(file, dest)
	if err != nil {
		return "", errors.New("fails to upload file")
	}

	hostTemplate, err := dal.GetFileServerPathTemplate()
	if err != nil {
		return "", errors.New("internal error")
	}

	// TODO random file name, and save file
	prefix := time.Now().Format("20060102")
	path := strings.ReplaceAll(hostTemplate, ":filename", prefix+file.Filename)
	return path, nil
}

func GetFilePath(userID string, filename string) (string, error) {
	// TODO valid user right

	// TODO convert file name to origin name

	return filename, nil
}
