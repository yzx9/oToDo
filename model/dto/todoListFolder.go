package dto

type TodoListMenu struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	IsLeaf   bool           `json:"isLeaf"`
	Children []TodoListMenu `json:"children"`
}
