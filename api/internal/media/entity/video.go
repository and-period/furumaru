package entity

import (
	"errors"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

var (
	ErrVideoRequiredProductIDs    = errors.New("entity: video required product ids")
	ErrVideoRequiredExperienceIDs = errors.New("entity: video required experience ids")
)

// VideoStatus - オンデマンド配信状況
type VideoStatus int32

const (
	VideoStatusUnknown   VideoStatus = 0
	VideoStatusPrivate   VideoStatus = 1 // 非公開
	VideoStatusWaiting   VideoStatus = 2 // 公開前
	VideoStatusLimited   VideoStatus = 3 // 限定公開
	VideoStatusPublished VideoStatus = 4 // 公開済み
)

// Video - オンデマンド配信情報
type Video struct {
	VideoProducts     `gorm:"-"`
	VideoExperiences  `gorm:"-"`
	ID                string      `gorm:"primaryKey;<-:create"` // オンデマンド動画ID
	CoordinatorID     string      `gorm:""`                     // コーディネータID
	ProductIDs        []string    `gorm:"-"`                    // 商品ID一覧
	ExperienceIDs     []string    `gorm:"-"`                    // 体験ID一覧
	Title             string      `gorm:""`                     // タイトル
	Description       string      `gorm:""`                     // 説明
	Status            VideoStatus `gorm:"-"`                    // 配信状況
	ThumbnailURL      string      `gorm:""`                     // サムネイルURL
	VideoURL          string      `gorm:""`                     // 動画URL
	Public            bool        `gorm:""`                     // 公開設定
	Limited           bool        `gorm:""`                     // 限定公開設定
	DisplayProduct    bool        `gorm:""`                     // 商品への表示設定
	DisplayExperience bool        `gorm:""`                     // 体験への表示設定
	PublishedAt       time.Time   `gorm:""`                     // 公開日時
	CreatedAt         time.Time   `gorm:"<-:create"`            // 作成日時
	UpdatedAt         time.Time   `gorm:""`                     // 更新日時
}

type Videos []*Video

type NewVideoParams struct {
	CoordinatorID     string
	ProductIDs        []string
	ExperienceIDs     []string
	Title             string
	Description       string
	ThumbnailURL      string
	VideoURL          string
	Public            bool
	Limited           bool
	DisplayProduct    bool
	DisplayExperience bool
	PublishedAt       time.Time
}

func NewVideo(params *NewVideoParams) *Video {
	videoID := uuid.Base58Encode(uuid.New())
	return &Video{
		ID:                videoID,
		CoordinatorID:     params.CoordinatorID,
		ProductIDs:        params.ProductIDs,
		ExperienceIDs:     params.ExperienceIDs,
		Title:             params.Title,
		Description:       params.Description,
		ThumbnailURL:      params.ThumbnailURL,
		VideoURL:          params.VideoURL,
		Public:            params.Public,
		Limited:           params.Limited,
		DisplayProduct:    params.DisplayProduct,
		DisplayExperience: params.DisplayExperience,
		PublishedAt:       params.PublishedAt,
		VideoProducts:     NewVideoProducts(videoID, params.ProductIDs),
		VideoExperiences:  NewVideoExperiences(videoID, params.ExperienceIDs),
	}
}

func (v *Video) Fill(products VideoProducts, experiences VideoExperiences, now time.Time) {
	v.ProductIDs = products.SortByPriority().ProductIDs()
	v.ExperienceIDs = experiences.SortByPriority().ExperienceIDs()
	v.VideoProducts = products
	v.VideoExperiences = experiences
	v.SetStatus(now)
}

func (v *Video) SetStatus(now time.Time) {
	switch {
	case !v.Public:
		v.Status = VideoStatusPrivate
	case now.Before(v.PublishedAt):
		v.Status = VideoStatusWaiting
	case v.Limited:
		v.Status = VideoStatusLimited
	default:
		v.Status = VideoStatusPublished
	}
}

func (v *Video) Published() bool {
	return v.Status == VideoStatusPublished || v.Status == VideoStatusLimited
}

func (vs Videos) IDs() []string {
	return set.UniqBy(vs, func(v *Video) string {
		return v.ID
	})
}

func (vs Videos) CoordinatorIDs() []string {
	return set.UniqBy(vs, func(v *Video) string {
		return v.CoordinatorID
	})
}

func (vs Videos) ProductIDs() []string {
	res := set.NewEmpty[string](len(vs))
	for i := range vs {
		res.Add(vs[i].ProductIDs...)
	}
	return res.Slice()
}

func (vs Videos) ExperienceIDs() []string {
	res := set.NewEmpty[string](len(vs))
	for i := range vs {
		res.Add(vs[i].ExperienceIDs...)
	}
	return res.Slice()
}

func (vs Videos) Fill(
	products map[string]VideoProducts,
	experiences map[string]VideoExperiences,
	now time.Time,
) {
	for i := range vs {
		vs[i].Fill(products[vs[i].ID], experiences[vs[i].ID], now)
	}
}
