package file

func (file File) CanAccessByUser(userID int64) bool {
	switch FileAccessType(file.AccessType) {
	case FileTypePublic:

	default:
		checker, ok := PermissionCheckerFactory.Get(file.AccessType)
		if !ok {
			return false
		}

		return checker(PermissionRequest{
			VisitorID: userID,
			FileID:    file.ID,
			RelatedID: file.RelatedID,
		})
	}

	return false
}

type PermissionRequest struct {
	VisitorID int64
	FileID    int64
	RelatedID int64
}

type permissionChecker = func(file PermissionRequest) bool

type permissionCheckerFactory struct {
	checkers map[FileAccessType]permissionChecker
}

var PermissionCheckerFactory permissionCheckerFactory = permissionCheckerFactory{
	checkers: make(map[FileAccessType]permissionChecker),
}

func (f *permissionCheckerFactory) Register(accessType FileAccessType, checker permissionChecker) {
	f.checkers[accessType] = checker
}

func (f *permissionCheckerFactory) Unregister(accessType FileAccessType) {
	delete(f.checkers, accessType)
}

func (f permissionCheckerFactory) Get(accessType FileAccessType) (permissionChecker, bool) {
	v, ok := f.checkers[accessType]
	return v, ok
}
