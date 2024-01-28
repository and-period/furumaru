package service

import "github.com/and-period/furumaru/api/internal/media/entity"

// UploadStatus - ファイルアップロード状況
type UploadStatus int32

const (
	UploadStatusUnknown   UploadStatus = 0
	UploadStatusWaiting   UploadStatus = 1 // アップロード待ち
	UploadStatusSucceeded UploadStatus = 2 // 成功
	UploadStatusFailed    UploadStatus = 3 // 失敗
)

func NewUploadStatus(status entity.UploadStatus) UploadStatus {
	switch status {
	case entity.UploadStatusSucceeded:
		return UploadStatusSucceeded
	case entity.UploadStatusFailed:
		return UploadStatusFailed
	case entity.UploadStatusWaiting:
		return UploadStatusWaiting
	default:
		return UploadStatusUnknown
	}
}

func (s UploadStatus) Response() int32 {
	return int32(s)
}
