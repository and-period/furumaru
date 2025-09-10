package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/set"
)

// お知らせ種別
type NotificationType types.NotificationType

// お知らせ通知先
type NotificationTarget types.NotificationTarget

type NotificationTargets []NotificationTarget

// お知らせ状態
type NotificationStatus types.NotificationStatus

type Notification struct {
	types.Notification
}

type Notifications []*Notification

func NewNotificationType(typ entity.NotificationType) NotificationType {
	switch typ {
	case entity.NotificationTypeSystem:
		return NotificationType(types.NotificationTypeSystem)
	case entity.NotificationTypeLive:
		return NotificationType(types.NotificationTypeLive)
	case entity.NotificationTypePromotion:
		return NotificationType(types.NotificationTypePromotion)
	case entity.NotificationTypeOther:
		return NotificationType(types.NotificationTypeOther)
	default:
		return NotificationType(types.NotificationTypeUnknown)
	}
}

func (t NotificationType) MessengerEntity() entity.NotificationType {
	switch types.NotificationType(t) {
	case types.NotificationTypeSystem:
		return entity.NotificationTypeSystem
	case types.NotificationTypeLive:
		return entity.NotificationTypeLive
	case types.NotificationTypePromotion:
		return entity.NotificationTypePromotion
	case types.NotificationTypeOther:
		return entity.NotificationTypeOther
	default:
		return entity.NotificationTypeUnknown
	}
}

func (t NotificationType) Response() types.NotificationType {
	return types.NotificationType(t)
}

func NewNotificationTarget(target entity.NotificationTarget) NotificationTarget {
	switch target {
	case entity.NotificationTargetAdministrators:
		return NotificationTarget(types.NotificationTargetAdministrators)
	case entity.NotificationTargetCoordinators:
		return NotificationTarget(types.NotificationTargetCoordinators)
	case entity.NotificationTargetProducers:
		return NotificationTarget(types.NotificationTargetProducers)
	case entity.NotificationTargetUsers:
		return NotificationTarget(types.NotificationTargetUsers)
	default:
		return NotificationTarget(types.NotificationTargetUnknown)
	}
}

func (t NotificationTarget) MessengerEntity() entity.NotificationTarget {
	switch types.NotificationTarget(t) {
	case types.NotificationTargetAdministrators:
		return entity.NotificationTargetAdministrators
	case types.NotificationTargetCoordinators:
		return entity.NotificationTargetCoordinators
	case types.NotificationTargetProducers:
		return entity.NotificationTargetProducers
	case types.NotificationTargetUsers:
		return entity.NotificationTargetUsers
	default:
		return entity.NotificationTargetUnknown
	}
}

func (t NotificationTarget) Response() types.NotificationTarget {
	return types.NotificationTarget(t)
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
		return NotificationStatus(types.NotificationStatusWaiting)
	case entity.NotificationStatusNotified:
		return NotificationStatus(types.NotificationStatusNotified)
	default:
		return NotificationStatus(types.NotificationStatusUnknown)
	}
}

func (s NotificationStatus) MessengerEntity() entity.NotificationStatus {
	switch types.NotificationStatus(s) {
	case types.NotificationStatusWaiting:
		return entity.NotificationStatusWaiting
	case types.NotificationStatusNotified:
		return entity.NotificationStatusNotified
	default:
		return entity.NotificationStatusUnknown
	}
}

func (s NotificationStatus) Response() types.NotificationStatus {
	return types.NotificationStatus(s)
}

func (ts NotificationTargets) Response() []types.NotificationTarget {
	res := make([]types.NotificationTarget, len(ts))
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
	if types.NotificationType(n.Type) == types.NotificationTypePromotion && promotion != nil {
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
