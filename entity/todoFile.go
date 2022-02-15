package entity

type TodoFile struct {
	Entity

	FileID   string `json:"fileID" gorm:"type:char(36);"`
	FileName string `json:"fileName" gorm:"-"`
	File     File   `json:"-"`

	TodoID string `json:"todoID" gorm:"type:char(36);"`
	Todo   Todo   `json:"-"`
}
