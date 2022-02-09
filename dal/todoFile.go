package dal

import (
	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
)

var todoFiles = make(map[uuid.UUID]entity.TodoFile)

func InsertTodoFile(file entity.TodoFile) error {
	todoFiles[file.ID] = file
	return nil
}

func GetTodoFiles(todoID uuid.UUID) ([]entity.TodoFile, error) {
	vec := make([]entity.TodoFile, 0)
	for _, v := range todoFiles {
		if v.TodoID == todoID {
			vec = append(vec, v)
		}
	}

	return vec, nil
}
