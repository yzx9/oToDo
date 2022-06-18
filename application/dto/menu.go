package dto

type MenuItem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`

	IsLeaf   bool       `json:"isLeaf"`
	Children []MenuItem `json:"children"`
}
