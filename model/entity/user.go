package entity

type User struct {
	Entity

	Name      string `json:"name" gorm:"size:128;uniqueIndex;"`
	Nickname  string `json:"nickname" gorm:"size:128"`
	Password  []byte `json:"-" gorm:"size:32;"`
	Email     string `json:"email" gorm:"size:32;"`
	Telephone string `json:"telephone" gorm:"size:16;"`
	Avatar    string `json:"avatar"`

	BasicTodoListID int64     `json:"basicTodoListID"`
	BasicTodoList   *TodoList `json:"-"`

	TodoLists []TodoList `json:"-"`

	SharedTodoLists []*TodoList `json:"-" gorm:"many2many:todo_list_shared_users"`
}