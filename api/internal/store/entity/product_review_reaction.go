package entity

import "time"

// ProductReviewReactionType - 商品レビューのリアクション種別
type ProductReviewReactionType int32

const (
	ProductReviewReactionTypeUnknown ProductReviewReactionType = 0
	ProductReviewReactionTypeLike    ProductReviewReactionType = 1 // いいね
	ProductReviewReactionTypeDislike ProductReviewReactionType = 2 // いまいち
)

// ProductReviewReaction - 商品レビューに対するリアクション
type ProductReviewReaction struct {
	ReviewID     string                    `gorm:"primaryKey;<-:create"` // 商品レビューID
	UserID       string                    `gorm:"primaryKey;<-:create"` // ユーザーID
	ReactionType ProductReviewReactionType `gorm:""`                     // リアクション種別
	CreatedAt    time.Time                 `gorm:"<-:create"`            // 作成日時
	UpdatedAt    time.Time                 `gorm:""`                     // 更新日時
}

type ProductReviewReactions []*ProductReviewReaction

type NewProductReviewReactionParams struct {
	ReviewID     string
	UserID       string
	ReactionType ProductReviewReactionType
}

func NewProductReviewReaction(params *NewProductReviewReactionParams) *ProductReviewReaction {
	return &ProductReviewReaction{
		ReviewID:     params.ReviewID,
		UserID:       params.UserID,
		ReactionType: params.ReactionType,
	}
}

// AggregatedProductReviewReaction - 商品レビューのリアクション
type AggregatedProductReviewReaction struct {
	ReviewID     string                    `gorm:""` // 商品レビューID
	ReactionType ProductReviewReactionType `gorm:""` // リアクション種別
	Total        int64                     `gorm:""` // リアクション数
}

type AggregatedProductReviewReactions []*AggregatedProductReviewReaction

func newEmptyProductReviewReactionTotal() map[ProductReviewReactionType]int64 {
	return map[ProductReviewReactionType]int64{
		ProductReviewReactionTypeLike:    0,
		ProductReviewReactionTypeDislike: 0,
	}
}

func (rs AggregatedProductReviewReactions) GetTotalByMap() map[ProductReviewReactionType]int64 {
	res := newEmptyProductReviewReactionTotal()
	for _, r := range rs {
		res[r.ReactionType] = r.Total
	}
	return res
}

func (rs AggregatedProductReviewReactions) GroupByReviewID() map[string]AggregatedProductReviewReactions {
	res := make(map[string]AggregatedProductReviewReactions, len(rs))
	for _, r := range rs {
		if _, ok := res[r.ReviewID]; !ok {
			res[r.ReviewID] = make(AggregatedProductReviewReactions, 0)
		}
		res[r.ReviewID] = append(res[r.ReviewID], r)
	}
	return res
}
