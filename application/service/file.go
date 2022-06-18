package service

import (
	"mime/multipart"
	"regexp"
	"strconv"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/util"
)

var supportedFileTypeRegex = regexp.MustCompile(`.(jpg|jpeg|JPG|png|PNG|gif|GIF|ico|ICO)$`)

func UploadPublicFile(f *multipart.FileHeader) (dto.File, error) {
	// only support img now
	if !supportedFileTypeRegex.MatchString(f.Filename) {
		return dto.File{}, util.NewErrorWithForbidden("unsupported file type")
	}

	record, err := file.UploadFile(file.FileTypePublic, 0, f)
	if err != nil {
		return dto.File{}, err
	}

	return dto.File{FileID: record.ID}, nil
}

func GetFilePath(userID int64, fileID string) (string, error) {
	id, err := strconv.ParseInt(fileID, 10, 64)

	var f *file.File
	if err == nil {
		f, err = file.GetFile(id)
	} else {
		f, err = file.GetFileByPreSignID(fileID)
	}

	if err != nil || !f.CanAccessByUser(userID) {
		return "", file.Notfound
	}

	return f.GetFilePath(), nil
}

func PreSignFile(payload dto.FilePreSign) (dto.FilePreSignResult, error) {
	file, err := file.GetFile(payload.FileID)
	if err != nil {
		return dto.FilePreSignResult{}, err
	}

	filePreSignID, err := file.CreateFilePreSignID(payload.UserID, payload.ExpiresIn)
	if err != nil {
		return dto.FilePreSignResult{}, err
	}

	return dto.FilePreSignResult{
		FileID: filePreSignID,
	}, nil
}
