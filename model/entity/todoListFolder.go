package entity

type TodoListFolder struct {
	Entity

	Name string `json:"name" gorm:"size:128"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TodoLists []TodoList `json:"-"`
}
