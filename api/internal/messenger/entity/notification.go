package entity

import (
	"encoding/json"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 掲載対象
type TargetType int32

const (
	PostTargetUnknown      TargetType = 0
	PostTargetUsers        TargetType = 1 // ユーザー対象
	PostTargetProducers    TargetType = 2 // 生産者対象
	PostTargetCoordinators TargetType = 3 // コーディネータ対象
)

var targetUsers = []int32{
	int32(PostTargetUsers),
}

var targetAdmins = []int32{
	int32(PostTargetCoordinators),
	int32(PostTargetProducers),
}

type NotificationOrderBy string

const (
	NotificationOrderByTitle       NotificationOrderBy = "title"
	NotificationOrderByPublic      NotificationOrderBy = "public"
	NotificationOrderByPublishedAt NotificationOrderBy = "published_at"
)

// Notification - お知らせ情報
type Notification struct {
	ID          string         `gorm:"primaryKey;<-:create"`        // お知らせID
	CreatedBy   string         `gorm:"<-:create"`                   // 登録者ID
	CreatorName string         `gorm:"<-:create"`                   // 登録者名
	UpdatedBy   string         `gorm:""`                            // 更新者ID
	Title       string         `gorm:""`                            // タイトル
	Body        string         `gorm:""`                            // 本文
	Targets     []TargetType   `gorm:"-"`                           // 掲載対象一覧
	TargetsJSON datatypes.JSON `gorm:"default:null;column:targets"` // 掲載対象一覧(JSON)
	PublishedAt time.Time      `gorm:""`                            // 掲載開始日時
	Public      bool           `gorm:""`                            // 公開フラグ
	CreatedAt   time.Time      `gorm:"<-:create"`                   // 作成日時
	UpdatedAt   time.Time      `gorm:""`                            // 更新日時
	DeletedAt   gorm.DeletedAt `gorm:"default:null"`                // 削除日時
}

type Notifications []*Notification

type NewNotificationParams struct {
	CreatedBy   string
	CreatorName string
	UpdatedBy   string
	Title       string
	Body        string
	Targets     []TargetType
	Public      bool
	PublishedAt time.Time
}

func NewNotification(params *NewNotificationParams) *Notification {
	return &Notification{
		ID:          uuid.Base58Encode(uuid.New()),
		CreatedBy:   params.CreatedBy,
		CreatorName: params.CreatorName,
		UpdatedBy:   params.UpdatedBy,
		Title:       params.Title,
		Body:        params.Body,
		Targets:     params.Targets,
		Public:      params.Public,
		PublishedAt: params.PublishedAt,
	}
}

func (n *Notification) HasUserTarget() bool {
	set := set.New(targetUsers...)
	for i := range n.Targets {
		if set.Contains(int32(n.Targets[i])) {
			return true
		}
	}
	return false
}

func (n *Notification) HasAdminTarget() bool {
	set := set.New(targetAdmins...)
	for i := range n.Targets {
		if set.Contains(int32(n.Targets[i])) {
			return true
		}
	}
	return false
}

func (n *Notification) HasCoordinatorTarget() bool {
	for i := range n.Targets {
		if n.Targets[i] == PostTargetCoordinators {
			return true
		}
	}
	return false
}

func (n *Notification) HasProducerTarget() bool {
	for i := range n.Targets {
		if n.Targets[i] == PostTargetProducers {
			return true
		}
	}
	return false
}

func (n *Notification) Fill() (err error) {
	n.Targets, err = n.unmarshalTarget()
	return
}

func (n *Notification) unmarshalTarget() ([]TargetType, error) {
	if n.TargetsJSON == nil {
		return []TargetType{}, nil
	}
	var targets []TargetType
	return targets, json.Unmarshal(n.TargetsJSON, &targets)
}

func (n *Notification) FillJSON() error {
	v, err := json.Marshal(n.Targets)
	if err != nil {
		return err
	}
	n.TargetsJSON = datatypes.JSON(v)
	return nil
}

func (ns Notifications) Fill() error {
	for i := range ns {
		if err := ns[i].Fill(); err != nil {
			return err
		}
	}
	return nil
}

func Marshal(t []TargetType) ([]byte, error) {
	if len(t) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(t)
}
