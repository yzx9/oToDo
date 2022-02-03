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

func UploadFile(file *multipart.FileHeader) (string, error) {
	if file.Size > maxFileSize {
		return "", utils.NewErrorWithHttpStatus(http.StatusRequestEntityTooLarge, "file too large")
	}

	record := entity.File{
		ID:                   uuid.New(),
		FileName:             file.Filename,
		CreatedAt:            time.Now(),
		FileDestTemplateID:   destTemplate.ID,
		FileServerTemplateID: serverTemplate.ID,
	}

	err := utils.SaveFile(file, applyTemplate(destTemplate.Template, record))
	if err != nil {
		return "", errors.New("fails to upload file")
	}

	err = dal.AddFile(record)
	if err != nil {
		return "", errors.New("fails to upload file")
	}

	path := applyTemplate(serverTemplate.Template, record)
	return path, nil
}

func GetFilePath(userID string, filename string) (string, error) {
	// Dest template is assumed to be foo/id.ext
	id := strings.TrimSuffix(filename, filepath.Ext(filename))
	uuid, err := uuid.Parse(id)
	if err != nil {
		return "", utils.NewErrorWithHttpStatus(http.StatusNotFound, "file not found: %v", filename)
	}

	file, err := dal.GetFile(uuid)
	if err != nil {
		return "", utils.NewErrorWithHttpStatus(http.StatusNotFound, "file not found: %v", filename)
	}

	// TODO valid user right

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
