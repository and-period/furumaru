package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
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
	HostPostalCode        string               `gorm:""`                                       // 開催場所(郵便番号)
	HostPrefecture        string               `gorm:"-"`                                      // 開催場所(都道府県)
	HostPrefectureCode    int32                `gorm:"column:host_prefecture"`                 // 開催場所(都道府県コード)
	HostCity              string               `gorm:""`                                       // 開催場所(市区町村)
	HostAddressLine1      string               `gorm:""`                                       // 開催場所(町名・番地)
	HostAddressLine2      string               `gorm:""`                                       // 開催場所(ビル名・号室など)
	HostGeolocation       mysql.Geometry       `gorm:""`                                       // 開催場所(座標情報)
	HostLongitude         float64              `gorm:"-"`                                      // 開催場所(座標情報:経度)
	HostLatitude          float64              `gorm:"-"`                                      // 開催場所(座標情報:緯度)
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

type NewExperienceParams struct {
	CoordinatorID         string
	ProducerID            string
	TypeID                string
	Title                 string
	Description           string
	Public                bool
	SoldOut               bool
	Media                 MultiExperienceMedia
	RecommendedPoints     []string
	PromotionVideoURL     string
	HostPostalCode        string
	HostPrefectureCode    int32
	HostCity              string
	HostAddressLine1      string
	HostAddressLine2      string
	HostLongitude         float64
	HostLatitude          float64
	StartAt               time.Time
	EndAt                 time.Time
	PriceAdult            int64
	PriceJuniorHighSchool int64
	PriceElementarySchool int64
	PricePreschool        int64
	PriceSenior           int64
}

func NewExperience(params *NewExperienceParams) (*Experience, error) {
	experienceID := uuid.Base58Encode(uuid.New())
	prefecture, err := codes.ToPrefectureJapanese(params.HostPrefectureCode)
	if err != nil {
		return nil, err
	}
	rparams := &NewExperienceRevisionParams{
		ExperienceID:          experienceID,
		PriceAdult:            params.PriceAdult,
		PriceJuniorHighSchool: params.PriceJuniorHighSchool,
		PriceElementarySchool: params.PriceElementarySchool,
		PricePreschool:        params.PricePreschool,
		PriceSenior:           params.PriceSenior,
	}
	revision := NewExperienceRevision(rparams)
	experience := &Experience{
		ID:                 experienceID,
		CoordinatorID:      params.CoordinatorID,
		ProducerID:         params.ProducerID,
		TypeID:             params.TypeID,
		Title:              params.Title,
		Description:        params.Description,
		Public:             params.Public,
		SoldOut:            params.SoldOut,
		Media:              params.Media,
		RecommendedPoints:  params.RecommendedPoints,
		PromotionVideoURL:  params.PromotionVideoURL,
		HostPostalCode:     params.HostPostalCode,
		HostPrefecture:     prefecture,
		HostPrefectureCode: params.HostPrefectureCode,
		HostCity:           params.HostCity,
		HostAddressLine1:   params.HostAddressLine1,
		HostAddressLine2:   params.HostAddressLine2,
		HostLongitude:      params.HostLongitude,
		HostLatitude:       params.HostLatitude,
		StartAt:            params.StartAt,
		EndAt:              params.EndAt,
		ExperienceRevision: *revision,
	}
	if err := experience.Validate(); err != nil {
		return nil, err
	}
	return experience, nil
}

func (e *Experience) Validate() error {
	if len(e.RecommendedPoints) > 3 {
		return errors.New("entity: limit exceeded recommended points")
	}
	if e.HostLongitude < -180 || 180 < e.HostLongitude {
		return errors.New("entity: invalid host longitude")
	}
	if e.HostLatitude < -90 || 90 < e.HostLatitude {
		return errors.New("entity: invalid host latitude")
	}
	return e.Media.Validate()
}

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
	e.SetLocation()
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

func (e *Experience) SetLocation() {
	e.HostLongitude = e.HostGeolocation.X
	e.HostLatitude = e.HostGeolocation.Y
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
	e.HostGeolocation = ExperienceHostGeolocation(e.HostLongitude, e.HostLatitude)
	return nil
}

func ExperienceMarshalRecommendedPoints(points []string) ([]byte, error) {
	return json.Marshal(points)
}

func ExperienceHostGeolocation(longitude, latitude float64) mysql.Geometry {
	return mysql.Geometry{X: longitude, Y: latitude}
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

func (es Experiences) Map() map[string]*Experience {
	res := make(map[string]*Experience, len(es))
	for _, e := range es {
		res[e.ID] = e
	}
	return res
}

func NewExperienceMedia(url string, isThumbnail bool) *ExperienceMedia {
	return &ExperienceMedia{
		URL:         url,
		IsThumbnail: isThumbnail,
	}
}

func (m MultiExperienceMedia) Validate() error {
	var exists bool
	for _, media := range m {
		if !media.IsThumbnail {
			continue
		}
		if exists {
			return errOnlyOneThumbnail
		}
		exists = true
	}
	return nil
}

func (m MultiExperienceMedia) Marshal() ([]byte, error) {
	if len(m) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(m)
}
