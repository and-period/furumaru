package entity

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ExperienceStatus - 体験販売状況
type ExperienceStatus int32

const (
	ExperienceStatusUnknown   ExperienceStatus = 0
	ExperienceStatusPrivate   ExperienceStatus = 1 // 非公開
	ExperienceStatusWaiting   ExperienceStatus = 2 // 公開前
	ExperienceStatusAccepting ExperienceStatus = 3 // 受付中
	ExperienceStatusSoldOut   ExperienceStatus = 4 // 定員オーバー
	ExperienceStatusFinished  ExperienceStatus = 5 // 終了済み
)

// Experience - 体験情報
type Experience struct {
	ExperienceRevision    `gorm:"-"`
	ID                    string               `gorm:"primaryKey;<-:create"`                   // 体験ID
	CoordinatorID         string               `gorm:""`                                       // コーディネータID
	ProducerID            string               `gorm:""`                                       // 生産者ID
	TypeID                string               `gorm:"column:experience_type_id"`              // 体験種別ID
	Title                 string               `gorm:""`                                       // タイトル
	Description           string               `gorm:""`                                       // 説明
	Public                bool                 `gorm:""`                                       // 公開フラグ
	SoldOut               bool                 `gorm:""`                                       // 定員オーバーフラグ
	Status                ExperienceStatus     `gorm:"-"`                                      // 販売状況
	ThumbnailURL          string               `gorm:"-"`                                      // サムネイルURL
	Media                 MultiExperienceMedia `gorm:"-"`                                      // メディア一覧
	MediaJSON             datatypes.JSON       `gorm:"default:null;column:media"`              // メディア一覧(JSON)
	RecommendedPoints     []string             `gorm:"-"`                                      // おすすめポイント一覧
	RecommendedPointsJSON datatypes.JSON       `gorm:"default:null;column:recommended_points"` // おすすめポイント一覧(JSON)
	PromotionVideoURL     string               `gorm:""`                                       // 紹介動画URL
	HostPrefecture        string               `gorm:"-"`                                      // 開催場所(都道府県)
	HostPrefectureCode    int32                `gorm:"column:origin_prefecture"`               // 開催場所(都道府県コード)
	HostCity              string               `gorm:""`                                       // 開催場所(市区町村)
	StartAt               time.Time            `gorm:""`                                       // 募集開始日時
	EndAt                 time.Time            `gorm:""`                                       // 募集終了日時
	CreatedAt             time.Time            `gorm:"<-:create"`                              // 登録日時
	UpdatedAt             time.Time            `gorm:""`                                       // 更新日時
	DeletedAt             gorm.DeletedAt       `gorm:"default:null"`                           // 削除日時
}

type Experiences []*Experience

// ExperienceMedia - 体験メディア情報
type ExperienceMedia struct {
	URL         string `json:"url"`         // メディアURL
	IsThumbnail bool   `json:"isThumbnail"` // サムネイルとして使用
}

type MultiExperienceMedia []*ExperienceMedia
