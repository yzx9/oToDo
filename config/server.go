package config

var Server = server{}

type server struct {
	ID                       string
	Port                     int
	Host                     string
	AccessControlAllowOrigin string
	FilePathTemplate         string // Support :id, :ext, :name, :path, :date
}
