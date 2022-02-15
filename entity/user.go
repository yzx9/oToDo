package entity

type User struct {
	Entity

	Name      string `json:"name" gorm:"size:128"`
	Nickname  string `json:"nickname" gorm:"size:128"`
	Password  []byte `json:"password" gorm:"size:32;"`
	Email     string `json:"email" gorm:"size:32;"`
	Telephone string `json:"telephone" gorm:"size:16;"`
	Avatar    string `json:"avatar"`

	BasicTodoListID string    `json:"basic_todo_list_id" gorm:"type:char(36);"`
	BasicTodoList   *TodoList `json:"-"`

	TodoLists []TodoList `json:"-"`
}
