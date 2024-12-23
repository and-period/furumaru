package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// ExperienceReview - 体験レビュー
type ExperienceReview struct {
	ID           string                                 `gorm:"primaryKey;<-:create"` // 阿智県レビューID
	ExperienceID string                                 `gorm:""`                     // 体験ID
	UserID       string                                 `gorm:""`                     // ユーザーID
	Rate         int64                                  `gorm:""`                     // 評価
	Title        string                                 `gorm:""`                     // タイトル
	Comment      string                                 `gorm:""`                     // コメント
	Reactions    map[ExperienceReviewReactionType]int64 `gorm:"-"`                    // リアクション
	CreatedAt    time.Time                              `gorm:"<-:create"`            // 作成日時
	UpdatedAt    time.Time                              `gorm:""`                     // 更新日時
	DeletedAt    gorm.DeletedAt                         `gorm:"default:null"`         // 削除日時
}

type ExperienceReviews []*ExperienceReview

// AggregatedExperienceReview - 体験レビュー集計情報
type AggregatedExperienceReview struct {
	ExperienceID string  `gorm:"primaryKey"` // 体験ID
	Count        int64   `gorm:""`           // レビュー数
	Average      float64 `gorm:""`           // 平均評価
	Rate1        int64   `gorm:""`           // 評価1の数
	Rate2        int64   `gorm:""`           // 評価2の数
	Rate3        int64   `gorm:""`           // 評価3の数
	Rate4        int64   `gorm:""`           // 評価4の数
	Rate5        int64   `gorm:""`           // 評価5の数
}

type AggregatedExperienceReviews []*AggregatedExperienceReview

type NewExperienceReviewParams struct {
	ExperienceID string
	UserID       string
	Rate         int64
	Title        string
	Comment      string
}

func NewExperienceReview(params *NewExperienceReviewParams) *ExperienceReview {
	return &ExperienceReview{
		ID:           uuid.Base58Encode(uuid.New()),
		ExperienceID: params.ExperienceID,
		UserID:       params.UserID,
		Rate:         params.Rate,
		Title:        params.Title,
		Comment:      params.Comment,
	}
}

func (r *ExperienceReview) SetReactions(reactions AggregatedExperienceReviewReactions) {
	r.Reactions = reactions.GetTotalByMap()
}

func (rs ExperienceReviews) SetReactions(reactions map[string]AggregatedExperienceReviewReactions) {
	for _, r := range rs {
		r.SetReactions(reactions[r.ID])
	}
}

func (rs ExperienceReviews) IDs() []string {
	return set.UniqBy(rs, func(r *ExperienceReview) string {
		return r.ID
	})
}

func (rs ExperienceReviews) UserIDs() []string {
	return set.UniqBy(rs, func(r *ExperienceReview) string {
		return r.UserID
	})
}

func (rs AggregatedExperienceReviews) Map() map[string]*AggregatedExperienceReview {
	m := make(map[string]*AggregatedExperienceReview)
	for _, r := range rs {
		m[r.ExperienceID] = r
	}
	return m
}
