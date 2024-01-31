package response

type UploadURLResponse struct {
	URL string `json:"url"` // アップロード用の署名付きURL
}

type UploadStateResponse struct {
	URL    string `json:"url"`    // 参照先URL
	Status int32  `json:"status"` // アップロード結果
}

type UploadImageResponse struct {
	URL string `json:"url"` // 画像アップロード先URL
}
