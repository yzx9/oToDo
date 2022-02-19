package entity

type SharingType = int8

const (
	SharingTypeTodoList SharingType = 10*iota + 1 // Set RelatedID to todo list id
)

type Sharing struct {
	Entity

	Token     string `json:"token" gorm:"size:128;uniqueIndex"`
	Active    bool   `json:"active"`
	Type      int8   `json:"type"`      // SharingType
	RelatedID int64  `json:"relatedID"` // Depends on Type

	UserID int64 `json:"userID"`
	User   User  `json:"-"`
}
