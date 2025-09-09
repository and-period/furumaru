package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/set"
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

type NotificationTargets []NotificationTarget

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

type Notification struct {
	types.Notification
}

type Notifications []*Notification

func NewNotificationType(typ entity.NotificationType) NotificationType {
	switch typ {
	case entity.NotificationTypeSystem:
		return NotificationTypeSystem
	case entity.NotificationTypeLive:
		return NotificationTypeLive
	case entity.NotificationTypePromotion:
		return NotificationTypePromotion
	case entity.NotificationTypeOther:
		return NotificationTypeOther
	default:
		return NotificationTypeUnknown
	}
}

func (t NotificationType) MessengerEntity() entity.NotificationType {
	switch t {
	case NotificationTypeSystem:
		return entity.NotificationTypeSystem
	case NotificationTypeLive:
		return entity.NotificationTypeLive
	case NotificationTypePromotion:
		return entity.NotificationTypePromotion
	case NotificationTypeOther:
		return entity.NotificationTypeOther
	default:
		return entity.NotificationTypeUnknown
	}
}

func (t NotificationType) Response() int32 {
	return int32(t)
}

func NewNotificationTarget(target entity.NotificationTarget) NotificationTarget {
	switch target {
	case entity.NotificationTargetAdministrators:
		return NotificationTargetAdministrators
	case entity.NotificationTargetCoordinators:
		return NotificationTargetCoordinators
	case entity.NotificationTargetProducers:
		return NotificationTargetProducers
	case entity.NotificationTargetUsers:
		return NotificationTargetUsers
	default:
		return NotificationTargetUnknown
	}
}

func (t NotificationTarget) MessengerEntity() entity.NotificationTarget {
	switch t {
	case NotificationTargetAdministrators:
		return entity.NotificationTargetAdministrators
	case NotificationTargetCoordinators:
		return entity.NotificationTargetCoordinators
	case NotificationTargetProducers:
		return entity.NotificationTargetProducers
	case NotificationTargetUsers:
		return entity.NotificationTargetUsers
	default:
		return entity.NotificationTargetUnknown
	}
}

func (t NotificationTarget) Response() int32 {
	return int32(t)
}

func NewNotificationTargets(targets []entity.NotificationTarget) NotificationTargets {
	res := make(NotificationTargets, len(targets))
	for i := range targets {
		res[i] = NewNotificationTarget(targets[i])
	}
	return res
}

func (ts NotificationTargets) MessengerEntities() []entity.NotificationTarget {
	res := make([]entity.NotificationTarget, len(ts))
	for i := range ts {
		res[i] = ts[i].MessengerEntity()
	}
	return res
}

func NewNotificationStatus(status entity.NotificationStatus) NotificationStatus {
	switch status {
	case entity.NotificationStatusWaiting:
		return NotificationStatusWaiting
	case entity.NotificationStatusNotified:
		return NotificationStatusNotified
	default:
		return NotificationStatusUnknown
	}
}

func (s NotificationStatus) MessengerEntity() entity.NotificationStatus {
	switch s {
	case NotificationStatusWaiting:
		return entity.NotificationStatusWaiting
	case NotificationStatusNotified:
		return entity.NotificationStatusNotified
	default:
		return entity.NotificationStatusUnknown
	}
}

func (s NotificationStatus) Response() int32 {
	return int32(s)
}

func (ts NotificationTargets) Response() []int32 {
	res := make([]int32, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}

func NewNotification(notification *entity.Notification) *Notification {
	return &Notification{
		Notification: types.Notification{
			ID:          notification.ID,
			Status:      NewNotificationStatus(notification.Status).Response(),
			Type:        NewNotificationType(notification.Type).Response(),
			Title:       notification.Title,
			Body:        notification.Body,
			Note:        notification.Note,
			Targets:     NewNotificationTargets(notification.Targets).Response(),
			PublishedAt: notification.PublishedAt.Unix(),
			PromotionID: notification.PromotionID,
			CreatedBy:   notification.CreatedBy,
			UpdatedBy:   notification.UpdatedBy,
			CreatedAt:   notification.CreatedAt.Unix(),
			UpdatedAt:   notification.UpdatedAt.Unix(),
		},
	}
}

func (n *Notification) Fill(promotion *Promotion) {
	if NotificationType(n.Type) == NotificationTypePromotion && promotion != nil {
		n.Title = promotion.Title
	}
}

func (n *Notification) Response() *types.Notification {
	return &n.Notification
}

func NewNotifications(notifications entity.Notifications) Notifications {
	res := make(Notifications, len(notifications))
	for i := range notifications {
		res[i] = NewNotification(notifications[i])
	}
	return res
}

func (ns Notifications) Response() []*types.Notification {
	res := make([]*types.Notification, len(ns))
	for i := range ns {
		res[i] = ns[i].Response()
	}
	return res
}

func (ns Notifications) AdminIDs() []string {
	return set.UniqBy(ns, func(n *Notification) string {
		return n.CreatedBy
	})
}

func (ns Notifications) Fill(promotions map[string]*Promotion) {
	for _, n := range ns {
		n.Fill(promotions[n.PromotionID])
	}
}
