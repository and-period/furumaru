package types

type GetUploadURLRequest struct {
	FileType string `json:"fileType" validate:"required"` // ファイル種別
}

type UploadURLResponse struct {
	Key string `json:"key"` // アップロード後の状態参照用キー
	URL string `json:"url"` // アップロード用の署名付きURL
}

type UploadStateResponse struct {
	URL    string `json:"url"`    // 参照先URL
	Status int32  `json:"status"` // アップロード結果
}
