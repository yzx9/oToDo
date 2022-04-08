package dto

type FileDTO struct {
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

type FilePreSignClaims struct {
	TokenClaims

	UserID int64 `json:"uid"`
	FileID int64 `json:"fileID"`
}
