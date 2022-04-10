package repository

import (
	"time"

	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type TodoStep struct {
	Entity

	Name   string
	Done   bool
	DoneAt *time.Time

	TodoID int64
	Todo   Todo
}

type TodoStepRepository struct {
	db *gorm.DB
}

func NewTodoStepRepository(db *gorm.DB) TodoStepRepository {
	return TodoStepRepository{db: db}
}

func (r TodoStepRepository) Save(entity *todo.TodoStep) error {
	po := r.convertToPO(entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
	return util.WrapGormErr(err, "todo step")
}

func (r TodoStepRepository) Delete(id int64) error {
	err := r.db.
		Delete(&TodoStep{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo step")
}

func (r TodoStepRepository) Find(id int64) (todo.TodoStep, error) {
	var entity TodoStep
	err := r.db.
		Where(&TodoStep{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&entity).
		Error

	return r.convertToEntity(entity), util.WrapGormErr(err, "todo step")
}

func (r TodoStepRepository) FindAllByTodo(todoID int64) ([]todo.TodoStep, error) {
	var POs []TodoStep
	err := r.db.
		Where(TodoStep{
			TodoID: todoID,
		}).
		Find(&POs).
		Error

	return r.convertToEntities(POs), util.WrapGormErr(err, "todo step")
}

func (r TodoStepRepository) convertToPO(entity *todo.TodoStep) TodoStep {
	return TodoStep{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Name:   entity.Name,
		Done:   entity.Done,
		DoneAt: entity.DoneAt,

		TodoID: entity.TodoID,
	}
}

func (r TodoStepRepository) convertToEntity(po TodoStep) todo.TodoStep {
	return todo.TodoStep{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		Name:   po.Name,
		Done:   po.Done,
		DoneAt: po.DoneAt,

		TodoID: po.TodoID,
	}
}

func (r TodoStepRepository) convertToEntities(POs []TodoStep) []todo.TodoStep {
	return util.Map(r.convertToEntity, POs)
}
