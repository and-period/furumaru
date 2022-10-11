package entity

import (
	"encoding/json"
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ライブ情報
type Live struct {
	ID             string         `gorm:"primaryKey;<-:create"`           // テンプレートID
	ScheduleID     string         `gorm:""`                               // 開催スケジュールID
	Title          string         `gorm:""`                               // タイトル
	Description    string         `gorm:""`                               // 説明
	ProducerID     string         `gorm:""`                               // 生産者ID
	StartAt        time.Time      `gorm:""`                               // 配信開始日時
	EndAt          time.Time      `gorm:""`                               // 配信終了日時
	Canceled       bool           `gorm:""`                               // 配信中止フラグ
	Recommends     []string       `gorm:"-"`                              // おすすめ商品一覧
	RecommendsJSON datatypes.JSON `gorm:"default:null;column:recommends"` // おすすめ商品一覧(JSON)
	CreatedAt      time.Time      `gorm:"<-:create"`                      // 登録日時
	UpdatedAt      time.Time      `gorm:""`                               // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`                   // 削除日時
}

type Lives []*Live

type NewLiveParams struct {
	ScheduleID  string
	Title       string
	Description string
	ProducerID  string
	StartAt     time.Time
	EndAt       time.Time
	Recommends  []string
}

func NewLive(params *NewLiveParams) *Live {
	return &Live{
		ID:          uuid.Base58Encode(uuid.New()),
		ScheduleID:  params.ScheduleID,
		Title:       params.Title,
		Description: params.Description,
		ProducerID:  params.ProducerID,
		StartAt:     params.StartAt,
		EndAt:       params.EndAt,
		Recommends:  params.Recommends,
	}
}

func (l *Live) FillJSON() error {
	v, err := Marshal(l.Recommends)
	if err != nil {
		return err
	}
	l.RecommendsJSON = datatypes.JSON(v)
	return nil
}

func Marshal(s []string) ([]byte, error) {
	if len(s) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(s)
}
