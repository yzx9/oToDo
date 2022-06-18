package service

import (
	"fmt"
	"mime/multipart"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/domain/todo"
)

func CreateTodo(userID int64, newTodo dto.NewTodo) (dto.Todo, error) {
	newTodo.UserID = userID // override user

	entity := newTodo.ToEntity()
	if err := entity.New(); err != nil {
		return dto.Todo{}, err
	}

	return dto.Todo{}.FromEntity(entity), nil
}

func UpdateTodo(userID int64, todo dto.Todo) (dto.Todo, error) {
	entity := todo.ToEntity()
	if err := entity.Save(userID); err != nil {
		return dto.Todo{}, nil
	}

	return dto.Todo{}.FromEntity(entity), nil
}

func DeleteTodo(userID int64, todoID int64) (dto.Todo, error) {
	t, err := TodoRepository.Find(todoID)
	if err != nil {
		return dto.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	if !t.CanAccessByUser(userID) {
		return dto.Todo{}, todo.PermissionDenied
	}

	if err := t.Delete(userID); err != nil {
		return dto.Todo{}, err
	}

	return dto.Todo{}.FromEntity(t), nil
}

func CreateTodoStep(userID int64, step dto.NewTodoStep) (dto.TodoStep, error) {
	t, err := TodoRepository.Find(step.TodoID)
	if err != nil {
		return dto.TodoStep{}, fmt.Errorf("fails to get todo: %w", err)
	}

	if !t.CanAccessByUser(userID) {
		return dto.TodoStep{}, todo.PermissionDenied
	}

	entity := t.NewStep()
	step.AssembleTo(&entity)
	if err := entity.New(); err != nil {
		return dto.TodoStep{}, err
	}

	return dto.TodoStep{}.FromEntity(entity), nil
}

func UpdateTodoStep(userID int64, step dto.TodoStep) (todo.TodoStep, error) {
	entity := step.ToEntity()
	if err := entity.Save(userID); err != nil {
		return todo.TodoStep{}, nil
	}

	return entity, nil
}

func DeleteTodoStep(userID, todoID, todoStepID int64) (todo.TodoStep, error) {
	t, err := TodoRepository.Find(todoID)
	if err != nil {
		return todo.TodoStep{}, fmt.Errorf("fails to get todo: %w", err)
	}

	if !t.CanAccessByUser(userID) {
		return todo.TodoStep{}, todo.PermissionDenied
	}

	step, err := t.GetStep(todoStepID)
	if err != nil {
		return todo.TodoStep{}, err
	}

	if err := step.Delete(); err != nil {
		return todo.TodoStep{}, err
	}

	return step, err
}

func UploadTodoFile(userID, todoID int64, file *multipart.FileHeader) (dto.File, error) {
	t, err := TodoRepository.Find(todoID)
	if err != nil {
		return dto.File{}, fmt.Errorf("fails to get todo: %w", err)
	}

	if !t.CanAccessByUser(userID) {
		return dto.File{}, todo.PermissionDenied
	}

	record, err := t.AddFile(file)
	if err != nil {
		return dto.File{}, err
	}

	return dto.File{FileID: record.ID}, nil
}

func CanAccessTodoFile(request file.PermissionRequest) bool {
	t, err := TodoRepository.Find(request.RelatedID)
	if err != nil {
		return false
	}

	return t.CanAccessByUser(request.VisitorID)
}

func GetTodo(userID, todoID int64) (todo.Todo, error) {
	t, err := TodoRepository.Find(todoID)
	if err != nil {
		return todo.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	if !t.CanAccessByUser(userID) {
		return todo.Todo{}, todo.PermissionDenied
	}

	return t, nil
}

func GetTodosByUserAndTodoList(userID, todoListID int64) ([]todo.Todo, error) {
	if _, err := todo.OwnOrSharedTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	todos, err := TodoRepository.FindAllByTodoList(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func GetTodosInBasicTodoList(userID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllInBasicTodoList(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func GetImportantTodosByUser(userID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllImportantOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get important todos: %w", err)
	}

	return todos, nil
}

func GetPlannedTodosByUser(userID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllPlanedOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get planed todos: %w", err)
	}

	return todos, nil
}

func GetNotNotifiedTodosByUser(userID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllNotNotifiedOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get not-notified todos: %w", err)
	}

	return todos, nil
}
