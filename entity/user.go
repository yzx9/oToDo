package entity

type User struct {
	Entity

	Name      string   `json:"name" gorm:"size:128"`
	Nickname  string   `json:"nickname" gorm:"size:128"`
	Password  [32]byte `json:"password"`
	Email     string   `json:"email"`
	Telephone string   `json:"telephone"`
	Avatar    string   `json:"avatar"`

	BasicTodoListID string    `json:"basic_todo_list_id" gorm:"type:char(36);"`
	BasicTodoList   *TodoList `json:"-"`
}
