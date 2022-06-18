package dto

type File struct {
	FileID int64 `json:"fileID"`
}

type FilePreSign struct {
	UserID    int64
	FileID    int64
	ExpiresIn int `json:"expiresIn"` // Unix
}

type FilePreSignResult struct {
	FileID string `json:"fileID"`
}
