package todo

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/infrastructure/util"
)

type Todo struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Title      string
	Memo       string
	Importance bool
	Deadline   *time.Time
	Notified   bool
	NotifyAt   *time.Time
	Done       bool
	DoneAt     *time.Time

	UserID int64

	TodoListID int64

	Files []int64

	Steps []int64

	TodoRepeatPlanID int64
	TodoRepeatPlan   TodoRepeatPlan

	NextID *int64 // next todo id if repeat
}

func GetTodoByUser(userID, todoID int64) (Todo, error) {
	todo, err := TodoRepository.Find(todoID)
	if err != nil {
		return Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	if _, err = todolist.OwnOrSharedTodoList(userID, todo.TodoListID); err != nil {
		return Todo{}, util.NewErrorWithForbidden("unable to handle non-owned todo: %v", todo.ID)
	}

	return todo, nil
}

func (todo *Todo) New() error {
	if _, err := todolist.OwnOrSharedTodoList(todo.UserID, todo.TodoListID); err != nil {
		return fmt.Errorf("fails to get todo list: %w", err)
	}

	plan, err := CreateTodoRepeatPlan(todo.TodoRepeatPlan)
	if err != nil {
		return fmt.Errorf("fails to create todo repeat plan: %w", err)
	}
	todo.TodoRepeatPlanID = plan.ID

	todo.ID = 0
	if err := TodoRepository.Save(todo); err != nil {
		return fmt.Errorf("fails to create todo: %w", err)
	}

	return nil
}

func (todo *Todo) Save(userID int64) error {
	// Limits
	oldTodo, err := GetTodoByUser(userID, todo.ID)
	if err != nil {
		return err
	}

	todo.CreatedAt = oldTodo.CreatedAt
	todo.UserID = oldTodo.UserID
	todo.Files = oldTodo.Files
	todo.Steps = oldTodo.Steps
	todo.NextID = oldTodo.NextID

	if !oldTodo.Done && todo.Done {
		t := time.Now()
		todo.DoneAt = &t

		// Create Repeat Todo If Need
		if todo.NextID == nil {
			created, next, err := CreateRepeatTodoIfNeed(*todo)
			if err != nil {
				return err
			}

			if created {
				todo.NextID = &next.ID
			}
		}
	}

	plan, err := UpdateTodoRepeatPlan(todo.TodoRepeatPlan, oldTodo.TodoRepeatPlan)
	if err != nil {
		return err
	}
	todo.TodoRepeatPlanID = plan.ID

	// Save
	if err = TodoRepository.Save(todo); err != nil {
		return err
	}

	go UpdateTagAsync(todo, oldTodo.Title)

	return nil
}

func (todo Todo) Delete(userID int64) error {
	if _, err := GetTodoByUser(userID, todo.ID); err != nil {
		return err
	}

	if err := TodoRepository.Delete(todo.ID); err != nil {
		return fmt.Errorf("fails to delete todo: %w", err)
	}

	go UpdateTagAsync(&todo, "")

	return nil
}

func (todo Todo) NewStep() TodoStep {
	return TodoStep{
		TodoID: todo.ID,
	}
}

func (todo Todo) GetStep(id int64) (TodoStep, error) {
	step, err := TodoStepRepository.Find(id)
	if err != nil {
		return TodoStep{}, fmt.Errorf("fails to get todo step: %w", err)
	}

	if step.TodoID != todo.ID {
		return TodoStep{}, fmt.Errorf("todo step not found")
	}

	return step, err
}
