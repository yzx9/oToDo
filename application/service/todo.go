package service

import (
	"fmt"
	"mime/multipart"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/domain/todolist"
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
	todo, err := todo.GetTodoByUser(userID, todoID)
	if err != nil {
		return dto.Todo{}, err
	}

	if err := todo.Delete(userID); err != nil {
		return dto.Todo{}, err
	}

	return dto.Todo{}.FromEntity(todo), nil
}

func CreateTodoStep(userID int64, step dto.NewTodoStep) (dto.TodoStep, error) {
	todo, err := todo.GetTodoByUser(userID, step.TodoID)
	if err != nil {
		return dto.TodoStep{}, nil
	}

	entity := todo.NewStep()
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
	entity, err := todo.GetTodoByUser(userID, todoID)
	if err != nil {
		return todo.TodoStep{}, err
	}

	step, err := entity.GetStep(todoStepID)
	if err != nil {
		return todo.TodoStep{}, err
	}

	if err := step.Delete(); err != nil {
		return todo.TodoStep{}, err
	}

	return step, err
}

func UploadTodoFile(userID, todoID int64, file *multipart.FileHeader) (dto.FileDTO, error) {
	todo, err := todo.GetTodoByUser(userID, todoID)
	if err != nil {
		return dto.FileDTO{}, err
	}

	record, err := todo.AddFile(file)
	if err != nil {
		return dto.FileDTO{}, err
	}

	return dto.FileDTO{FileID: record.ID}, nil
}

func CanAccessTodoFile(request file.PermissionRequest) bool {
	_, err := todo.GetTodoByUser(request.VisitorID, request.RelatedID)
	return err == nil
}

func GetTodo(userID, todoID int64) (todo.Todo, error) {
	entity, err := todo.GetTodoByUser(userID, todoID)
	if err != nil {
		return todo.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	return entity, nil
}

func GetTodosByTodoList(todoListID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllByTodoList(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func GetTodosByUserAndTodoList(userID, todoListID int64) ([]todo.Todo, error) {
	if _, err := todolist.OwnOrSharedTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	return GetTodosByTodoList(todoListID)
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
