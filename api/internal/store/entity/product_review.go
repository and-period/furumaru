package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// ProductReview - 商品レビュー
type ProductReview struct {
	ID        string         `gorm:"primaryKey;<-:create"` // 商品レビューID
	ProductID string         `gorm:""`                     // 商品ID
	UserID    string         `gorm:""`                     // ユーザーID
	Rate      int64          `gorm:""`                     // 評価
	Title     string         `gorm:""`                     // タイトル
	Comment   string         `gorm:""`                     // コメント
	CreatedAt time.Time      `gorm:"<-:create"`            // 作成日時
	UpdatedAt time.Time      `gorm:""`                     // 更新日時
	DeletedAt gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type ProductReviews []*ProductReview

// AggregatedProductReview - 商品レビュー集計情報
type AggregatedProductReview struct {
	ProductID string  `gorm:"primaryKey"` // 商品ID
	Count     int64   `gorm:""`           // レビュー数
	Average   float64 `gorm:""`           // 平均評価
	Rate1     int64   `gorm:""`           // 評価1の数
	Rate2     int64   `gorm:""`           // 評価2の数
	Rate3     int64   `gorm:""`           // 評価3の数
	Rate4     int64   `gorm:""`           // 評価4の数
	Rate5     int64   `gorm:""`           // 評価5の数
}

type AggregatedProductReviews []*AggregatedProductReview

type NewProductReviewParams struct {
	ProductID string
	UserID    string
	Rate      int64
	Title     string
	Comment   string
}

func NewProductReview(params *NewProductReviewParams) *ProductReview {
	return &ProductReview{
		ID:        uuid.Base58Encode(uuid.New()),
		ProductID: params.ProductID,
		UserID:    params.UserID,
		Rate:      params.Rate,
		Title:     params.Title,
		Comment:   params.Comment,
	}
}

func (rs ProductReviews) UserIDs() []string {
	return set.UniqBy(rs, func(r *ProductReview) string {
		return r.UserID
	})
}

func (rs AggregatedProductReviews) Map() map[string]*AggregatedProductReview {
	m := make(map[string]*AggregatedProductReview)
	for _, r := range rs {
		m[r.ProductID] = r
	}
	return m
}
