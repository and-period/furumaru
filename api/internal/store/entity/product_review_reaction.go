package entity

import "time"

type ProductReviewReactionType int32

const (
	ProductReviewReactionTypeUnknown ProductReviewReactionType = 0
	ProductReviewReactionTypeLike    ProductReviewReactionType = 1 // いいね
	ProductReviewReactionTypeDislike ProductReviewReactionType = 2 // いまいち
)

type ProductReviewReaction struct {
	ReviewID  string                    `gorm:"primaryKey;<-:create"` // 商品レビューID
	UserID    string                    `gorm:"primaryKey;<-:create"` // ユーザーID
	Type      ProductReviewReactionType `gorm:""`                     // リアクション種別
	CreatedAt time.Time                 `gorm:""`                     // 作成日時
	UpdatedAt time.Time                 `gorm:""`                     // 更新日時
}

type ProductReviewReactions []*ProductReviewReaction

type NewProductReviewReactionParams struct {
	ReviewID string
	UserID   string
	Type     ProductReviewReactionType
}

func NewProductReviewReaction(params *NewProductReviewReactionParams) *ProductReviewReaction {
	return &ProductReviewReaction{
		ReviewID: params.ReviewID,
		UserID:   params.UserID,
		Type:     params.Type,
	}
}
