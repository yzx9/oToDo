package config

var Database = database{}

type database struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	Protocol     string
	DatabaseName string
}
