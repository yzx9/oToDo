package sharing

/**
 * Repository
 */

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
