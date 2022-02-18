package entity

type SharingType = string

const (
	SharingTypeTodoList SharingType = "todoList"
)

type Sharing struct {
	Entity

	Token   string `json:"token" gorm:"size:128;uniqueIndex"`
	Active  bool   `json:"active"`
	Type    string `json:"type" gorm:"size:10"` // SharingType
	Payload string `json:"payload"`

	UserID string `json:"userID" gorm:"type:char(36);"`
	User   User   `json:"-"`
}
