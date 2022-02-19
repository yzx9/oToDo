package dto

type TodoListMenuItemRaw struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Count            int    `json:"count"`
	TodoListFolderID int64  `json:"-"`
}

type TodoListMenuItem struct {
	TodoListMenuItemRaw

	IsLeaf   bool               `json:"isLeaf"`
	Children []TodoListMenuItem `json:"children"`
}
