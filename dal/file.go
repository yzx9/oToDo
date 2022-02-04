package dal

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
)

var filePathTemplates = make(map[uuid.UUID]entity.FilePathTemplate)
var files = make(map[uuid.UUID]entity.File)

func InsertFile(file entity.File) error {
	files[file.ID] = file
	return nil
}

func GetFile(id uuid.UUID) (entity.File, error) {
	file, ok := files[id]
	if !ok {
		return entity.File{}, errors.New("file not found")
	}

	return file, nil
}

// FilePathTemplate

func GetFilePathTemplates() ([]entity.FilePathTemplate, error) {
	vec := []entity.FilePathTemplate{}
	for _, template := range filePathTemplates {
		vec = append(vec, template)
	}
	return vec, nil
}
