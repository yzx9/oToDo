package entity

type SharingType = int8

const (
	SharingTypeTodoList SharingType = 10*iota + 1 // Set RelatedID to todo list id
)

type Sharing struct {
	Entity

	Token     string `json:"-" gorm:"size:128;uniqueIndex"`
	Active    bool   `json:"active"`
	Type      int8   `json:"type"` // SharingType
	RelatedID int64  `json:"-"`    // Depends on Type

	UserID int64 `json:"-"`
	User   User  `json:"-"`
}
