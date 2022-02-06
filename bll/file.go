package bll

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

const maxFileSize = 8 << 20 // 8MiB

var destTemplate entity.FilePathTemplate
var serverTemplate entity.FilePathTemplate

func init() {
	// TODO db init
	var err error
	destTemplate, err = switchFilePathTemplate(entity.FilePathTemplateTypeDest)
	if err != nil {
		panic(err)
	}

	serverTemplate, err = switchFilePathTemplate(entity.FilePathTemplateTypeServer)
	if err != nil {
		panic(err)
	}
}

func UploadTodoFile(todoID string, file *multipart.FileHeader) (string, error) {
	id, err := uuid.Parse(todoID)
	if err != nil {
		return "", fmt.Errorf("invalid todo id: %v", todoID)
	}

	fileID := uuid.New()
	path, err := uploadFile(file, entity.File{
		ID:                   fileID,
		FileName:             file.Filename,
		FilePath:             "", // TODO add path template
		AccessType:           string(entity.FileTypeTodo),
		RelatedID:            id,
		CreatedAt:            time.Now(),
		FileDestTemplateID:   destTemplate.ID,
		FileServerTemplateID: serverTemplate.ID,
	})
	if err != nil {
		return "", fmt.Errorf("fails to upload file, %w", err)
	}

	_, err = dal.InsertTodoFile(id, entity.TodoFile{
		ID:     uuid.New(),
		FileID: fileID,
	})
	if err != nil {
		return "", fmt.Errorf("fails to upload todo file, %w", err)
	}

	return path, err
}

func uploadFile(file *multipart.FileHeader, record entity.File) (string, error) {
	if file.Size > maxFileSize {
		return "", utils.NewErrorWithHttpStatus(http.StatusRequestEntityTooLarge, "file too large")
	}

	err := utils.SaveFile(file, applyTemplate(destTemplate.Template, record))
	if err != nil {
		return "", errors.New("fails to upload file")
	}

	err = dal.InsertFile(record)
	if err != nil {
		return "", errors.New("fails to upload file")
	}

	return record.ID.String(), nil
}

func GetFile(fileID string) (entity.File, error) {
	uuid, err := uuid.Parse(fileID)
	if err != nil {
		return entity.File{}, utils.NewErrorWithNotFound("file not found, file id: %v", fileID)
	}

	file, err := dal.GetFile(uuid)
	if err != nil {
		return entity.File{}, utils.NewErrorWithNotFound("file not found, file id: %v", fileID)
	}

	return file, nil
}

func GetFileServerPath(fileID string) (string, error) {
	file, err := GetFile(fileID)
	if err != nil {
		return "", err
	}

	path := applyTemplate(serverTemplate.Template, file)
	return path, nil
}

func GetFilePath(fileID string) (string, error) {
	file, err := GetFile(fileID)
	if err != nil {
		return "", err
	}

	path := applyTemplate(destTemplate.Template, file)
	return path, nil
}

func GetFilePathWithAuth(filename string, userID string) (string, error) {
	// Dest template is assumed to be foo/fileID.ext
	fileID := strings.TrimSuffix(filename, filepath.Ext(filename))
	file, err := GetFile(fileID)
	if err != nil {
		return "", err
	}

	err = hasPermission(&file, userID)
	if err != nil {
		return "", err
	}

	path := applyTemplate(destTemplate.Template, file)
	return path, nil
}

func switchFilePathTemplate(templateType entity.FilePathTemplateType) (entity.FilePathTemplate, error) {
	var rawType string
	switch templateType {
	case entity.FilePathTemplateTypeDest:
		rawType = string(entity.FilePathTemplateTypeDest)

	case entity.FilePathTemplateTypeServer:
		rawType = string(entity.FilePathTemplateTypeServer)

	default:
		return entity.FilePathTemplate{}, fmt.Errorf("invalid template type")
	}

	templates, err := dal.GetFilePathTemplates()
	if err != nil {
		return entity.FilePathTemplate{}, fmt.Errorf("fails to get templates, %v", err)
	}

	for _, template := range templates {
		if template.Available && template.Type == rawType {
			return template, nil
		}
	}

	return entity.FilePathTemplate{}, fmt.Errorf("template not assign, type: %v", rawType)
}

func applyTemplate(template string, file entity.File) string {
	template = strings.ReplaceAll(template, ":date", file.CreatedAt.Format("20060102"))
	template = strings.ReplaceAll(template, ":filename", file.ID.String())
	template = strings.ReplaceAll(template, ":ext", filepath.Ext(file.FileName))
	return template
}

func hasPermission(file *entity.File, userID string) error {
	switch entity.FileAccessType(file.AccessType) {
	case entity.FileTypePublic:
		return nil

	case entity.FileTypeTodo:
		user, err := dal.GetUserByTodo(file.RelatedID)
		if err != nil {
			return fmt.Errorf("fails to get user, %w", err)
		}

		if userID != user.ID.String() {
			return utils.NewErrorWithHttpStatus(http.StatusUnauthorized, "no permission")
		}

	default:
		return fmt.Errorf("invalid file access type: %v", file.AccessType)
	}

	return nil
}
