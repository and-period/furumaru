package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/media/entity"
)

// UploadStatus - ファイルアップロード状況
type UploadStatus types.UploadStatus

func NewUploadStatus(status entity.UploadStatus) UploadStatus {
	switch status {
	case entity.UploadStatusSucceeded:
		return UploadStatus(types.UploadStatusSucceeded)
	case entity.UploadStatusFailed:
		return UploadStatus(types.UploadStatusFailed)
	case entity.UploadStatusWaiting:
		return UploadStatus(types.UploadStatusWaiting)
	default:
		return UploadStatus(types.UploadStatusUnknown)
	}
}

func (s UploadStatus) Response() types.UploadStatus {
	return types.UploadStatus(s)
}
