package entity

import (
	"errors"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotificationAlreadyPublished = errors.New(
		"entity: this notification is already published",
	)
	ErrNotificationDuplicatedTargets   = errors.New("entity: duplicated notification targets")
	ErrNotificationIncorrectTargets    = errors.New("entity: incorrect notification targets")
	ErrNotificationRequiredPromotionID = errors.New("entity: required promotion id")
	ErrNotificationRequiredTitle       = errors.New("entity: required title")
)

// お知らせ種別
type NotificationType int32

const (
	NotificationTypeUnknown   NotificationType = 0
	NotificationTypeOther     NotificationType = 1 // その他
	NotificationTypeSystem    NotificationType = 2 // システム関連
	NotificationTypeLive      NotificationType = 3 // ライブ関連
	NotificationTypePromotion NotificationType = 4 // セール関連
)

// お知らせ通知先
type NotificationTarget int32

const (
	NotificationTargetUnknown        NotificationTarget = 0
	NotificationTargetUsers          NotificationTarget = 1 // ユーザー
	NotificationTargetProducers      NotificationTarget = 2 // 生産者
	NotificationTargetCoordinators   NotificationTarget = 3 // コーディネータ
	NotificationTargetAdministrators NotificationTarget = 4 // 管理者
)

// お知らせ状態
type NotificationStatus int32

const (
	NotificationStatusUnknown  NotificationStatus = 0
	NotificationStatusWaiting  NotificationStatus = 1 // 投稿前
	NotificationStatusNotified NotificationStatus = 2 // 投稿済み
)

var targetUsers = []int32{
	int32(NotificationTargetUsers),
}

var targetAdmins = []int32{
	int32(NotificationTargetProducers),
	int32(NotificationTargetCoordinators),
	int32(NotificationTargetAdministrators),
}

// Notification - お知らせ情報
type Notification struct {
	ID          string               `gorm:"primaryKey;<-:create"` // お知らせID
	Type        NotificationType     `gorm:""`                     // お知らせ種別
	Status      NotificationStatus   `gorm:"-"`                    // お知らせ状態
	Title       string               `gorm:""`                     // タイトル
	Body        string               `gorm:""`                     // 本文
	Note        string               `gorm:""`                     // 備考
	Targets     []NotificationTarget `gorm:"-"`                    // お知らせ通知先一覧
	PublishedAt time.Time            `gorm:""`                     // 掲載開始日時
	PromotionID string               `gorm:"default:null"`         // プロモーションID
	CreatedBy   string               `gorm:"<-:create"`            // 登録者ID
	UpdatedBy   string               `gorm:""`                     // 更新者ID
	CreatedAt   time.Time            `gorm:"<-:create"`            // 作成日時
	UpdatedAt   time.Time            `gorm:""`                     // 更新日時
	DeletedAt   gorm.DeletedAt       `gorm:"default:null"`         // 削除日時
}

type Notifications []*Notification

type NewNotificationParams struct {
	Type        NotificationType
	Targets     []NotificationTarget
	Title       string
	Body        string
	Note        string
	PublishedAt time.Time
	PromotionID string
	CreatedBy   string
}

func NewNotification(params *NewNotificationParams) *Notification {
	return &Notification{
		ID:          uuid.Base58Encode(uuid.New()),
		Type:        params.Type,
		Targets:     params.Targets,
		Title:       params.Title,
		Body:        params.Body,
		Note:        params.Note,
		PublishedAt: params.PublishedAt,
		PromotionID: params.PromotionID,
		CreatedBy:   params.CreatedBy,
		UpdatedBy:   params.CreatedBy,
	}
}

func (n *Notification) TemplateID() MessageTemplateID {
	switch n.Type {
	case NotificationTypeSystem:
		return MessageTemplateIDNotificationSystem
	case NotificationTypeLive:
		return MessageTemplateIDNotificationLive
	case NotificationTypePromotion:
		return MessageTemplateIDNotificationPromotion
	default:
		return MessageTemplateIDNotificationOther
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

func (n *Notification) HasAdministratorTarget() bool {
	for i := range n.Targets {
		if n.Targets[i] == NotificationTargetAdministrators {
			return true
		}
	}
	return false
}

func (n *Notification) HasCoordinatorTarget() bool {
	for i := range n.Targets {
		if n.Targets[i] == NotificationTargetCoordinators {
			return true
		}
	}
	return false
}

func (n *Notification) HasProducerTarget() bool {
	for i := range n.Targets {
		if n.Targets[i] == NotificationTargetProducers {
			return true
		}
	}
	return false
}

func (n *Notification) Validate(now time.Time) error {
	if now.After(n.PublishedAt) {
		return ErrNotificationAlreadyPublished
	}
	if len(n.Targets) < 1 || len(n.Targets) > 4 {
		return ErrNotificationIncorrectTargets
	}
	targets := set.Uniq(n.Targets...)
	if len(targets) != len(n.Targets) {
		return ErrNotificationDuplicatedTargets
	}
	switch n.Type {
	case NotificationTypePromotion:
		if n.PromotionID == "" {
			return ErrNotificationRequiredPromotionID
		}
	default:
		if n.Title == "" {
			return ErrNotificationRequiredTitle
		}
	}
	return nil
}

func (n *Notification) Fill(now time.Time) {
	n.SetStatus(now)
}

func (n *Notification) SetStatus(now time.Time) {
	if now.After(n.PublishedAt) {
		n.Status = NotificationStatusNotified
	} else {
		n.Status = NotificationStatusWaiting
	}
}

func (ns Notifications) PromotionIDs() []string {
	set := set.NewEmpty[string](len(ns))
	for i := range ns {
		if ns[i].Type != NotificationTypePromotion {
			continue
		}
		set.Add(ns[i].PromotionID)
	}
	return set.Slice()
}

func (ns Notifications) AdminIDs() []string {
	set := set.NewEmpty[string](len(ns) * 2)
	for i := range ns {
		set.Add(ns[i].CreatedBy, ns[i].UpdatedBy)
	}
	return set.Slice()
}

func (ns Notifications) Fill(now time.Time) {
	for i := range ns {
		ns[i].Fill(now)
	}
}
