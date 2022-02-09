package bll

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

const maxFileSize = 8 << 20                        // 8MiB
const fileDestTemplate = "tmp/files/:date/:id:ext" // TODO Configurable

func UploadTodoFile(todoID uuid.UUID, file *multipart.FileHeader) (uuid.UUID, error) {
	fileID := uuid.New()
	path, err := uploadFile(file, entity.File{
		ID:         fileID,
		FileName:   file.Filename,
		AccessType: string(entity.FileTypeTodo),
		RelatedID:  todoID,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	err = dal.InsertTodoFile(entity.TodoFile{
		ID:     uuid.New(),
		FileID: fileID,
		TodoID: todoID,
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("fails to upload todo file: %w", err)
	}

	return path, err
}

func uploadFile(file *multipart.FileHeader, record entity.File) (uuid.UUID, error) {
	if file.Size > maxFileSize {
		return uuid.UUID{}, utils.NewErrorWithHttpStatus(http.StatusRequestEntityTooLarge, "file too large")
	}

	record.FilePath = applyTemplate(fileDestTemplate, record)
	err := utils.SaveFile(file, record.FilePath)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("fails to upload file: %w", err)
	}

	err = dal.InsertFile(record)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("fails to upload file: %w", err)
	}

	return record.ID, nil
}

func GetFile(fileID uuid.UUID) (entity.File, error) {
	file, err := dal.GetFile(fileID)
	if err != nil {
		return entity.File{}, utils.NewErrorWithNotFound("file not found: %v", fileID)
	}

	return file, nil
}

func GetFilePath(userID, fileID uuid.UUID) (string, error) {
	file, err := OwnFile(userID, fileID)
	if err != nil {
		return "", err
	}

	path := applyTemplate(fileDestTemplate, file)
	return path, nil
}

func applyTemplate(template string, file entity.File) string {
	template = strings.ReplaceAll(template, ":id", file.ID.String())
	template = strings.ReplaceAll(template, ":ext", filepath.Ext(file.FileName))
	template = strings.ReplaceAll(template, ":name", file.FileName)
	template = strings.ReplaceAll(template, ":path", file.FilePath)
	template = strings.ReplaceAll(template, ":date", file.CreatedAt.Format("20060102"))
	return template
}

func OwnFile(userID, fileID uuid.UUID) (entity.File, error) {
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
