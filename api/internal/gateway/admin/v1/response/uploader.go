package response

type UploadURLResponse struct {
	URL string `json:"url"` // アップロード用の署名付きURL
}

type UploadImageResponse struct {
	URL string `json:"url"` // 画像アップロード先URL
}
