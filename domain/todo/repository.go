package todo

import "github.com/yzx9/otodo/infrastructure/repository"

var TodoRepository todoRepository
var TodoStepRepository todoStepRepository
var TodoRepeatPlanRepository todoRepeatPlanRepository
var TagRepository tagRepository
var TagTodoRepository tagTodoRepository

type todoRepository interface {
	Save(todo *repository.Todo) error

	Delete(id int64) error

	Find(id int64) (repository.Todo, error)
}

type todoStepRepository interface {
	Save(todoStep *repository.TodoStep) error

	Delete(id int64) error

	Find(id int64) (repository.TodoStep, error)
}

type todoRepeatPlanRepository interface {
	Save(plan *repository.TodoRepeatPlan) error

	Delete(id int64) error

	Find(id int64) (repository.TodoRepeatPlan, error)
}

type tagRepository interface {
	Save(tag *repository.Tag) error

	Exist(userID int64, tagName string) (bool, error)
}

type tagTodoRepository interface {
	Save(userID, todoID int64, tagName string) error

	Delete(userID, todoID int64, tagName string) error
}
