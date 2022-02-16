package dto

type FileDTO struct {
	FileID int64 `json:"fileID"`
}

type FilePreSignDTO struct {
	ExpiresIn int `json:"expiresIn"` // Unix
}

type FilePreSignResultDTO struct {
	FileID string `json:"fileID"`
}

type FilePreSignClaims struct {
	TokenClaims

	UserID int64 `json:"uid"`
	FileID int64 `json:"fileID"`
}
