package entity

import (
	"encoding/json"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/set"
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
	ExperienceStatusArchived  ExperienceStatus = 6 // アーカイブ済み
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
	HostPrefectureCode    int32                `gorm:"column:host_prefecture"`                 // 開催場所(都道府県コード)
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

func (e *Experience) Fill(revision *ExperienceRevision, now time.Time) (err error) {
	e.Media, err = e.unmarshalMedia()
	if err != nil {
		return
	}
	e.RecommendedPoints, err = e.unmarshalRecommendedPoints()
	if err != nil {
		return
	}
	e.SetStatus(now)
	e.SetThumbnail()
	e.ExperienceRevision = *revision
	e.HostPrefecture, _ = codes.ToPrefectureJapanese(e.HostPrefectureCode)
	return
}

func (e *Experience) SetStatus(now time.Time) {
	switch {
	case !e.DeletedAt.Time.IsZero():
		e.Status = ExperienceStatusArchived
	case !e.Public:
		e.Status = ExperienceStatusPrivate
	case e.SoldOut:
		e.Status = ExperienceStatusSoldOut
	case now.Before(e.StartAt):
		e.Status = ExperienceStatusWaiting
	case now.Before(e.EndAt):
		e.Status = ExperienceStatusAccepting
	default:
		e.Status = ExperienceStatusFinished
	}
}

func (e *Experience) SetThumbnail() {
	for _, media := range e.Media {
		if !media.IsThumbnail {
			continue
		}
		e.ThumbnailURL = media.URL
	}
}

func (e *Experience) unmarshalMedia() (MultiExperienceMedia, error) {
	if e.MediaJSON == nil {
		return MultiExperienceMedia{}, nil
	}
	var media MultiExperienceMedia
	return media, json.Unmarshal(e.MediaJSON, &media)
}

func (e *Experience) unmarshalRecommendedPoints() ([]string, error) {
	if e.RecommendedPointsJSON == nil {
		return []string{}, nil
	}
	var points []string
	return points, json.Unmarshal(e.RecommendedPointsJSON, &points)
}

func (e *Experience) FillJSON() error {
	media, err := e.Media.Marshal()
	if err != nil {
		return err
	}
	points, err := ExperienceMarshalRecommendedPoints(e.RecommendedPoints)
	if err != nil {
		return err
	}
	e.MediaJSON = media
	e.RecommendedPointsJSON = points
	return nil
}

func ExperienceMarshalRecommendedPoints(points []string) ([]byte, error) {
	return json.Marshal(points)
}

func (es Experiences) Fill(revisions map[string]*ExperienceRevision, now time.Time) error {
	for _, e := range es {
		revision, ok := revisions[e.ID]
		if !ok {
			revision = &ExperienceRevision{ExperienceID: e.ID}
		}
		if err := e.Fill(revision, now); err != nil {
			return err
		}
	}
	return nil
}

func (es Experiences) IDs() []string {
	return set.UniqBy(es, func(e *Experience) string {
		return e.ID
	})
}

func (es Experiences) CoordinatorIDs() []string {
	return set.UniqBy(es, func(e *Experience) string {
		return e.CoordinatorID
	})
}

func (es Experiences) ProducerIDs() []string {
	return set.UniqBy(es, func(e *Experience) string {
		return e.ProducerID
	})
}

func (es Experiences) ExperienceTypeIDs() []string {
	return set.UniqBy(es, func(e *Experience) string {
		return e.TypeID
	})
}

func (m MultiExperienceMedia) Marshal() ([]byte, error) {
	if len(m) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(m)
}
