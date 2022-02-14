package entity

type TodoFile struct {
	Entity

	FileID   string `json:"file_id" gorm:"size:36"`
	FileName string `json:"file_name" gorm:"-"`
	File     File   `json:"-"`

	TodoID string `json:"todo_id" gorm:"size:36"`
	Todo   Todo   `json:"-"`
}
