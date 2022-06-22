package entity

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// 掲載対象
type TargetType int32

const (
	PostTargetAll          TargetType = 0
	PostTargetUsers        TargetType = 1
	PostTargetProducers    TargetType = 2
	PostTargetCoordinators TargetType = 3
)

type PostTarget struct {
	PostTarget TargetType `json:"PostTarget"`
}

type PostTargetList []*PostTarget

// Notification - お知らせ情報
type Notification struct {
	ID          string         `gorm:"primaryKey;<-:create"`        // お知らせID
	CreatedBy   string         `gorm:""`                            // 登録者ID
	CreatorName string         `gorm:""`                            // 登録者名
	UpdatedBy   string         `gorm:""`                            // 更新者ID
	Title       string         `gorm:""`                            // タイトル
	Body        string         `gorm:""`                            // 本文
	PublishedAt time.Time      `gorm:""`                            // 掲載開始日時
	Targets     PostTargetList `gorm:"-"`                           // 掲載対象一覧
	TargetsJSON datatypes.JSON `gorm:"default:null;column:targets"` // 掲載対象一覧(JSON)
}

func (n *Notification) Fill() error {
	var targets PostTargetList
	if err := json.Unmarshal(n.TargetsJSON, &targets); err != nil {
		return err
	}
	n.Targets = targets
	return nil
}

func (n *Notification) FillJSON() error {
	v, err := json.Marshal(n.Targets)
	if err != nil {
		return err
	}
	n.TargetsJSON = datatypes.JSON(v)
	return nil
}
