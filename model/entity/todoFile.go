package entity

type TodoFile struct {
	Entity

	FileID   int64  `json:"fileID"`
	FileName string `json:"fileName" gorm:"-"`
	File     File   `json:"-"`

	TodoID int64 `json:"todoID"`
	Todo   Todo  `json:"-"`
}
