package bll

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
)

const maxFileSize = 8 << 20 // 8MiB

var supportedFileTypeRegex = regexp.MustCompile(`.(jpg|jpeg|JPG|png|PNG|gif|GIF|ico|ICO)$`)

func UploadPublicFile(file *multipart.FileHeader) (int64, error) {
	// only support img now
	if !supportedFileTypeRegex.MatchString(file.Filename) {
		return 0, util.NewError(otodo.ErrPreconditionFailed, "invalid file type")
	}

	record := entity.File{
		FileName:   file.Filename,
		AccessType: int8(entity.FileTypePublic),
	}
	err := uploadFile(file, &record)
	return record.ID, err
}

func UploadTodoFile(todoID int64, file *multipart.FileHeader) (int64, error) {
	record := entity.File{
		FileName:   file.Filename,
		AccessType: int8(entity.FileTypeTodo),
		RelatedID:  todoID,
	}
	if err := uploadFile(file, &record); err != nil {
		return 0, err
	}

	if err := dal.InsertTodoFile(&entity.TodoFile{
		File: record,
		Todo: entity.Todo{Entity: entity.Entity{ID: record.ID}},
	}); err != nil {
		return 0, fmt.Errorf("fails to upload todo file: %w", err)
	}

	return record.ID, nil
}

func uploadFile(file *multipart.FileHeader, record *entity.File) error {
	write := func(err error) error {
		return fmt.Errorf("fails to upload file: %w", err)
	}

	if file.Size > maxFileSize {
		return util.NewError(otodo.ErrRequestEntityTooLarge, "file too large")
	}

	if err := dal.InsertFile(record); err != nil {
		return write(err)
	}

	// TODO[pref]: avoid duplicate save, remove :id in template?
	record.FileServerID = otodo.Conf.Server.ID
	record.FilePath = applyFilePathTemplate(record)
	if err := util.SaveFile(file, record.FilePath); err != nil {
		return write(err)
	}

	if err := dal.SaveFile(record); err != nil {
		return write(err)
	}

	return nil
}

func GetFile(fileID int64) (*entity.File, error) {
	file, err := dal.SelectFile(fileID)
	return file, fmt.Errorf("fails to get file: %w", err)
}

func GetFilePath(userID, fileID int64) (string, error) {
	file, err := OwnFile(userID, fileID)
	if err != nil {
		return "", err
	}

	return getFilePath(file), nil
}

func ForceGetFilePath(fileID int64) (string, error) {
	file, err := GetFile(fileID)
	if err != nil {
		return "", err
	}

	return getFilePath(file), nil
}

func OwnFile(userID, fileID int64) (*entity.File, error) {
	write := func(err error) (*entity.File, error) {
		return nil, err
	}

	file, err := GetFile(fileID)
	if err != nil {
		return write(err)
	}

	switch entity.FileAccessType(file.AccessType) {
	case entity.FileTypePublic:
		break

	case entity.FileTypeTodo:
		// TODO should check right from todo list
		user, err := dal.SelectUserByTodo(file.RelatedID)
		if err != nil {
			return write(fmt.Errorf("fails to get user: %w", err))
		}

		if userID != user.ID {
			return write(util.NewErrorWithForbidden("unable to get non-owned file: %v", file.ID))
		}

	default:
		return write(fmt.Errorf("invalid file access type: %v", file.AccessType))
	}

	return file, nil
}

func applyFilePathTemplate(file *entity.File) string {
	template := otodo.Conf.Server.FilePathTemplate
	template = strings.ReplaceAll(template, ":id", strconv.FormatInt(int64(file.ID), 10))
	template = strings.ReplaceAll(template, ":ext", filepath.Ext(file.FileName))
	template = strings.ReplaceAll(template, ":name", file.FileName)
	template = strings.ReplaceAll(template, ":path", file.FilePath)
	template = strings.ReplaceAll(template, ":date", file.CreatedAt.Format("20060102"))
	return template
}

func getFilePath(file *entity.File) string {
	// TODO[feat]: If exist multi servers, how to get file? maybe we need redirect
	return file.FilePath
}
