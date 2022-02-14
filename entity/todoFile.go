package entity

type TodoFile struct {
	ID       string `json:"id"`
	FileID   string `json:"file_id"`
	FileName string `json:"file_name" gorm:"-"`
	File     File   `json:"-"`

	TodoID string `json:"todo_id"`
	Todo   Todo   `json:"-"`
}
