package request

type GetUploadURLRequest struct {
	FileType string `json:"fileType" validate:"required"` // ファイル種別
}
