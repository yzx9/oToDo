package dto

type TodoListMenuItem struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Count            int    `json:"count"`
	TodoListFolderID int64  `json:"-"`

	IsLeaf   bool               `json:"isLeaf"`
	Children []TodoListMenuItem `json:"children"`
}
