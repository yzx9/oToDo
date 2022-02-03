package dal

// Replace :filename with real file name
func GetFileDestTemplate() (string, error) {
	// TODO Configurable
	return "./file/:filename", nil
}

// Replace :filename with real file name
func GetFileServerPathTemplate() (string, error) {
	// TODO Configurable
	return "http://localhost:8080/file/:filename", nil
}
