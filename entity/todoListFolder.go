package entity

type TodoListFolder struct {
	Entity

	Name string `json:"name"`

	UserID string `json:"user_id" gorm:"size:36"`
	User   User   `json:"-"`

	TodoLists []TodoList `json:"-"`
}
