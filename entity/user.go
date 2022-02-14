package entity

type User struct {
	Entity

	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	Password  []byte `json:"password"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	Avatar    string `json:"avatar"`

	BasicTodoListID string    `json:"basic_todo_list_id" gorm:"size:36"`
	BasicTodoList   *TodoList `json:"-"`
}
