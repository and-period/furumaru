package request

type GetUploadURLRequest struct {
	FileType string `json:"fileType,omitempty"` // ファイル種別
}
