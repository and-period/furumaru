package types

// UploadStatus - ファイルアップロード状況
type UploadStatus int32

const (
	UploadStatusUnknown   UploadStatus = 0
	UploadStatusWaiting   UploadStatus = 1 // アップロード待ち
	UploadStatusSucceeded UploadStatus = 2 // 成功
	UploadStatusFailed    UploadStatus = 3 // 失敗
)

type GetUploadURLRequest struct {
	FileType string `json:"fileType" validate:"required"` // ファイル種別
}

type UploadURLResponse struct {
	Key string `json:"key"` // アップロード後の状態参照用キー
	URL string `json:"url"` // アップロード用の署名付きURL
}

type UploadStateResponse struct {
	URL    string       `json:"url"`    // 参照先URL
	Status UploadStatus `json:"status"` // アップロード結果
}
