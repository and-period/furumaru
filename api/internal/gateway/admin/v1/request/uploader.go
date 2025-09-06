package request

type GetUploadURLRequest struct {
	FileType string `json:"fileType" binding:"required"` // ファイル種別
}
