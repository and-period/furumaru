package entity

import (
	"time"

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
	CreatedAt time.Time      `gorm:""`                     // 作成日時
	UpdatedAt time.Time      `gorm:""`                     // 更新日時
	DeletedAt gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type ProductReviews []*ProductReview

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
