package todolist

type SharingType = int8

const (
	SharingTypeTodoList SharingType = 10*iota + 1 // Set RelatedID to todo list id
)
