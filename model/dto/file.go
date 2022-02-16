package dto

type FileDTO struct {
	FileID string `json:"fileID"`
}

type FilePreSignDTO struct {
	ExpiresIn int `json:"expiresIn"` // Unix
}

type FilePreSignClaims struct {
	TokenClaims

	UserID string `json:"uid"`
	FileID string `json:"fileID"`
}
