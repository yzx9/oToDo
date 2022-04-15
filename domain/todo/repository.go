package todo

var TodoRepository todoRepository
var TodoStepRepository todoStepRepository
var TodoRepeatPlanRepository todoRepeatPlanRepository
var TagRepository tagRepository
var TagTodoRepository tagTodoRepository
var TodoFileRepository todoFileRepository // TODO move to todo domain

type todoRepository interface {
	Save(todo *Todo) error

	Delete(id int64) error

	Find(id int64) (Todo, error)
}

type todoStepRepository interface {
	Save(entity *TodoStep) error

	Delete(id int64) error

	Find(id int64) (TodoStep, error)
}

type todoRepeatPlanRepository interface {
	Save(entity *TodoRepeatPlan) error

	Delete(id int64) error

	Find(id int64) (TodoRepeatPlan, error)
}

type tagRepository interface {
	Save(tag *Tag) error

	Exist(userID int64, tagName string) (bool, error)
}

type tagTodoRepository interface {
	Save(userID, todoID int64, tagName string) error

	Delete(userID, todoID int64, tagName string) error
}

type todoFileRepository interface {
	Save(todoID, fileID int64) error
}
