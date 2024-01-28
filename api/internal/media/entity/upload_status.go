package entity

import "time"

// UploadStatus - ファイルアップロード状況
type UploadStatus int32

const (
	UploadStatusUnknown   UploadStatus = 0
	UploadStatusWaiting   UploadStatus = 1 // アップロード待ち
	UploadStatusSucceeded UploadStatus = 2 // 成功
	UploadStatusFailed    UploadStatus = 3 // 失敗
)

// UploadEvent - ファイルアップロード履歴情報
type UploadEvent struct {
	Key          string       `dynamodbav:"key"`                 // オブジェクトキー
	Status       UploadStatus `dynamodbav:"status"`              // アップロード状況
	FileGroup    string       `dynamodbav:"file_group"`          // ファイル種別（アップロード先ディレクトリ）
	FileType     string       `dynamodbav:"file_type"`           // MINEタイプ
	UploadURL    string       `dynamodbav:"upload_url"`          // ファイルアップロード先URL（一時保管用）
	ReferenceURL string       `dynamodbav:"reference_url"`       // ファイルアップロード先URL（参照用）
	ExpiredAt    time.Time    `dynamodbav:"expired_at,unixtime"` // 有効期限
	CreatedAt    time.Time    `dynamodbav:"created_at"`          // 登録日時
	UpdatedAt    time.Time    `dynamodbav:"updated_at"`          // 更新日時
}

type UploadEventParams struct {
	Key       string
	FileGroup string
	FileType  string
	UploadURL string
	Now       time.Time
	TTL       time.Duration
}

func NewUploadEvent(params *UploadEventParams) *UploadEvent {
	return &UploadEvent{
		Key:          params.Key,
		Status:       UploadStatusWaiting,
		FileGroup:    params.FileGroup,
		FileType:     params.FileType,
		UploadURL:    params.UploadURL,
		ReferenceURL: "",
		ExpiredAt:    params.Now.Add(params.TTL),
		CreatedAt:    params.Now,
		UpdatedAt:    params.Now,
	}
}

func (e *UploadEvent) TableName() string {
	return "upload-events"
}

func (e *UploadEvent) PrimaryKey() map[string]interface{} {
	return map[string]interface{}{
		"key": e.Key,
	}
}

func (e *UploadEvent) SetResult(success bool, referenceURL string, now time.Time) {
	if success {
		e.Status = UploadStatusSucceeded
		e.ReferenceURL = referenceURL
		e.UpdatedAt = now
	} else {
		e.Status = UploadStatusFailed
		e.UpdatedAt = now
	}
}

func (e *UploadEvent) Reguration() (*Regulation, error) {
	switch e.FileGroup {
	// ライブ配信関連
	case BroadcastArchiveMP4Path:
		return BroadcastArchiveRegulation, nil
	case BroadcastLiveMP4Path:
		return BroadcastLiveMP4Regulation, nil
	// コーディネータ関連
	case CoordinatorThumbnailPath:
		return CoordinatorThumbnailRegulation, nil
	case CoordinatorHeaderPath:
		return CoordinatorHeaderRegulation, nil
	case CoordinatorPromotionVideoPath:
		return CoordinatorPromotionVideoRegulation, nil
	case CoordinatorBonusVideoPath:
		return CoordinatorBonusVideoRegulation, nil
	// 生産者関連
	case ProducerThumbnailPath:
		return ProducerThumbnailRegulation, nil
	case ProducerHeaderPath:
		return ProducerHeaderRegulation, nil
	case ProducerPromotionVideoPath:
		return ProducerPromotionVideoRegulation, nil
	case ProducerBonusVideoPath:
		return ProducerBonusVideoRegulation, nil
	// 購入者関連
	case UserThumbnailPath:
		return UserThumbnailRegulation, nil
	// 商品関連
	case ProductMediaImagePath:
		return ProductMediaImageRegulation, nil
	case ProductMediaVideoPath:
		return ProductMediaVideoRegulation, nil
	// 品目関連
	case ProductTypeIconPath:
		return ProductTypeIconRegulation, nil
	// 開催スケジュール関連
	case ScheduleThumbnailPath:
		return ScheduleThumbnailRegulation, nil
	case ScheduleImagePath:
		return ScheduleImageRegulation, nil
	case ScheduleOpeningVideoPath:
		return ScheduleOpeningVideoRegulation, nil
	default:
		return nil, ErrNotFoundReguration
	}
}
