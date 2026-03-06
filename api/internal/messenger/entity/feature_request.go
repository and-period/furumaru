package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// 要望リクエストステータス
type FeatureRequestStatus int32

const (
	FeatureRequestStatusWaiting    FeatureRequestStatus = iota + 1 // 受付中
	FeatureRequestStatusReviewing                                   // 検討中
	FeatureRequestStatusAdopted                                     // 採用
	FeatureRequestStatusInProgress                                  // 対応中
	FeatureRequestStatusDone                                        // 完了
	FeatureRequestStatusRejected                                    // 却下
)

// 要望リクエストカテゴリ
type FeatureRequestCategory int32

const (
	FeatureRequestCategoryUI          FeatureRequestCategory = iota + 1 // UI
	FeatureRequestCategoryFeature                                        // 機能
	FeatureRequestCategoryPerformance                                    // パフォーマンス
	FeatureRequestCategoryOther                                          // その他
)

// 要望リクエスト優先度
type FeatureRequestPriority int32

const (
	FeatureRequestPriorityLow    FeatureRequestPriority = iota + 1 // 低
	FeatureRequestPriorityMedium                                    // 中
	FeatureRequestPriorityHigh                                      // 高
)

// FeatureRequest - 要望リクエスト情報
type FeatureRequest struct {
	ID          string                 `gorm:"primaryKey;<-:create"` // 要望リクエストID
	Title       string                 `gorm:""`                     // 件名
	Description string                 `gorm:""`                     // 内容
	Category    FeatureRequestCategory `gorm:""`                     // カテゴリ
	Priority    FeatureRequestPriority `gorm:""`                     // 優先度
	Status      FeatureRequestStatus   `gorm:""`                     // ステータス
	Note        string                 `gorm:""`                     // 管理者コメント
	SubmittedBy string                 `gorm:""`                     // 提出者 admin ID
	CreatedAt   time.Time              `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time              `gorm:""`                     // 更新日時
	DeletedAt   gorm.DeletedAt         `gorm:"default:null"`         // 削除日時
}

type FeatureRequests []*FeatureRequest

type NewFeatureRequestParams struct {
	Title       string
	Description string
	Category    FeatureRequestCategory
	Priority    FeatureRequestPriority
	SubmittedBy string
}

func NewFeatureRequest(params *NewFeatureRequestParams) *FeatureRequest {
	return &FeatureRequest{
		ID:          uuid.Base58Encode(uuid.New()),
		Title:       params.Title,
		Description: params.Description,
		Category:    params.Category,
		Priority:    params.Priority,
		Status:      FeatureRequestStatusWaiting,
		Note:        "",
		SubmittedBy: params.SubmittedBy,
	}
}
