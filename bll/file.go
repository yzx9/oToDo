package bll

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/utils"
)

const maxFileSize = 8 << 20 // 8MiB

func UploadTodoFile(todoID string, file *multipart.FileHeader) (string, error) {
	fileID := uuid.NewString()
	path, err := uploadFile(file, entity.File{
		Entity: entity.Entity{
			ID:        fileID,
			CreatedAt: time.Now(),
		},
		FileName:   file.Filename,
		AccessType: string(entity.FileTypeTodo),
		RelatedID:  todoID,
	})
	if err != nil {
		return "", err
	}

	err = dal.InsertTodoFile(entity.TodoFile{
		Entity: entity.Entity{
			ID: uuid.NewString(),
		},
		FileID: fileID,
		TodoID: todoID,
	})
	if err != nil {
		return "", fmt.Errorf("fails to upload todo file: %w", err)
	}

	return path, err
}

func uploadFile(file *multipart.FileHeader, record entity.File) (string, error) {
	if file.Size > maxFileSize {
		return "", utils.NewError(otodo.ErrRequestEntityTooLarge, "file too large")
	}

	record.FileServerID = otodo.Conf.Server.ID
	record.FilePath = applyFilePathTemplate(record)
	err := utils.SaveFile(file, record.FilePath)
	if err != nil {
		return "", fmt.Errorf("fails to upload file: %w", err)
	}

	err = dal.InsertFile(record)
	if err != nil {
		return "", fmt.Errorf("fails to upload file: %w", err)
	}

	return record.ID, nil
}

func GetFile(fileID string) (entity.File, error) {
	file, err := dal.GetFile(fileID)
	if err != nil {
		return entity.File{}, utils.NewErrorWithNotFound("file not found: %v", fileID)
	}

	return file, nil
}

func GetFilePath(userID, fileID string) (string, error) {
	file, err := OwnFile(userID, fileID)
	if err != nil {
		return "", err
	}

	path := applyFilePathTemplate(file)
	return path, nil
}

func ForceGetFilePath(fileID string) (string, error) {
	file, err := GetFile(fileID)
	if err != nil {
		return "", err
	}

	path := applyFilePathTemplate(file)
	return path, nil
}

func applyFilePathTemplate(file entity.File) string {
	template := otodo.Conf.Server.FilePathTemplate
	template = strings.ReplaceAll(template, ":id", file.ID)
	template = strings.ReplaceAll(template, ":ext", filepath.Ext(file.FileName))
	template = strings.ReplaceAll(template, ":name", file.FileName)
	template = strings.ReplaceAll(template, ":path", file.FilePath)
	template = strings.ReplaceAll(template, ":date", file.CreatedAt.Format("20060102"))
	return template
}

func OwnFile(userID, fileID string) (entity.File, error) {
	r := func(err error) (entity.File, error) {
		return entity.File{}, err
	}

	file, err := GetFile(fileID)
	if err != nil {
		return r(err)
	}

	switch entity.FileAccessType(file.AccessType) {
	case entity.FileTypePublic:
		break

	case entity.FileTypeTodo:
		user, err := dal.GetUserByTodo(file.RelatedID)
		if err != nil {
			return r(fmt.Errorf("fails to get user: %w", err))
		}

		if userID != user.ID {
			return r(utils.NewErrorWithForbidden("unable to get non-owned file: %v", file.ID))
		}

	default:
		return r(fmt.Errorf("invalid file access type: %v", file.AccessType))
	}

	return file, nil
}
