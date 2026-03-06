package types

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
	ID            string                 `json:"id"`            // 要望リクエストID
	Title         string                 `json:"title"`         // 件名
	Description   string                 `json:"description"`   // 内容
	Category      FeatureRequestCategory `json:"category"`      // カテゴリ
	Priority      FeatureRequestPriority `json:"priority"`      // 優先度
	Status        FeatureRequestStatus   `json:"status"`        // ステータス
	Note          string                 `json:"note"`          // 管理者コメント
	SubmittedBy   string                 `json:"submittedBy"`   // 提出者 admin ID
	SubmitterName string                 `json:"submitterName"` // 提出者名
	CreatedAt     int64                  `json:"createdAt"`     // 登録日時
	UpdatedAt     int64                  `json:"updatedAt"`     // 更新日時
}

type CreateFeatureRequestRequest struct {
	Title       string                 `json:"title" validate:"required,max=128"`       // 件名
	Description string                 `json:"description" validate:"required,max=2000"` // 内容
	Category    FeatureRequestCategory `json:"category" validate:"required"`             // カテゴリ
	Priority    FeatureRequestPriority `json:"priority" validate:"required"`             // 優先度
}

type UpdateFeatureRequestRequest struct {
	Status FeatureRequestStatus `json:"status" validate:"required"` // ステータス
	Note   string               `json:"note" validate:"max=2000"`   // 管理者コメント
}

type FeatureRequestResponse struct {
	FeatureRequest *FeatureRequest `json:"featureRequest"` // 要望リクエスト情報
}

type FeatureRequestsResponse struct {
	FeatureRequests []*FeatureRequest `json:"featureRequests"` // 要望リクエスト一覧
	Total           int64             `json:"total"`           // 合計件数
}
