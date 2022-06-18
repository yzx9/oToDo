package todo

/**
 * Event publisher
 */

var EventPublisher interface {
	Publish(event string, payload []byte)
	Subscribe(event string, cb func([]byte)) func()
}

/**
 * Repository
 */

var TodoRepository interface {
	Save(todo *Todo) error
	Delete(id int64) error
	Find(id int64) (Todo, error)
	DeleteAllByTodoList(todoListID int64) (int64, error)
}

var TodoStepRepository interface {
	Save(entity *TodoStep) error
	Delete(id int64) error
	Find(id int64) (TodoStep, error)
}

var TodoRepeatPlanRepository interface {
	Save(entity *TodoRepeatPlan) error
	Delete(id int64) error
	Find(id int64) (TodoRepeatPlan, error)
}

var TagRepository interface {
	Save(tag *Tag) error
	Exist(userID int64, tagName string) (bool, error)
}

var TagTodoRepository interface {
	Save(userID, todoID int64, tagName string) error
	Delete(userID, todoID int64, tagName string) error
}

var TodoFileRepository interface {
	Save(todoID, fileID int64) error
}

var TodoListRepository interface {
	Save(entity *TodoList) error
	Delete(id int64) error
	DeleteAllByFolder(todoListFolderID int64) (int64, error)
	Find(id int64) (TodoList, error)
	FindByUserOnMenuFormat(userID int64) ([]MenuItem, error)
}

var TodoListFolderRepository interface {
	Save(entity *TodoListFolder) error
	Delete(id int64) error
	Find(id int64) (TodoListFolder, error)
	FindAllByUser(userId int64) ([]TodoListFolder, error)
}

var TodoListSharingRepository interface {
	SaveSharedUser(userID, todoListID int64) error
	DeleteSharedUser(userID, todoListID int64) error
	ExistSharing(userID, todoListID int64) (bool, error)
}
