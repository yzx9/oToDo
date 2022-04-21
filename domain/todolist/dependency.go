package todolist

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
	DeleteAllByTodoList(todoListID int64) (int64, error)
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

var SharingRepository interface {
	Save(entity *Sharing) error
	DeleteAllByUserAndType(userID int64, sharingType SharingType) (int64, error)
	Find(token string) (Sharing, error)
}

var TodoListSharingRepository interface {
	SaveSharedUser(userID, todoListID int64) error
	DeleteSharedUser(userID, todoListID int64) error
	ExistSharing(userID, todoListID int64) (bool, error)
}
