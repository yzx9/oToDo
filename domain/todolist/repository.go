package todolist

var TodoRepository todoRepository
var TodoListRepository todoListRepository
var TodoListFolderRepository todoListFolderRepository
var SharingRepository sharingRepository
var TodoListSharingRepository todoListSharingRepository

type todoRepository interface {
	DeleteAllByTodoList(todoListID int64) (int64, error)
}

type todoListRepository interface {
	Save(entity *TodoList) error

	Delete(id int64) error

	DeleteAllByFolder(todoListFolderID int64) (int64, error)

	Find(id int64) (TodoList, error)
}

type todoListFolderRepository interface {
	Save(entity *TodoListFolder) error

	Delete(id int64) error

	Find(id int64) (TodoListFolder, error)
}

type sharingRepository interface {
	Save(entity *Sharing) error

	DeleteAllByUserAndType(userID int64, sharingType SharingType) (int64, error)

	Find(token string) (Sharing, error)
}

type todoListSharingRepository interface {
	SaveSharedUser(userID, todoListID int64) error

	DeleteSharedUser(userID, todoListID int64) error

	ExistSharing(userID, todoListID int64) (bool, error)
}
