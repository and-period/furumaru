package service

import (
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	set "github.com/and-period/furumaru/api/pkg/set/v2"
)

type Notification struct {
	response.Notification
}

type Notifications []*Notification

type TargetType struct {
	response.TargetType
}

type TargetTypes []*TargetType

func NewNotification(notification *entity.Notification) *Notification {
	return &Notification{
		Notification: response.Notification{
			ID:          notification.ID,
			CreatedBy:   notification.CreatedBy,
			CreatorName: notification.CreatorName,
			UpdatedBy:   notification.UpdatedBy,
			Title:       notification.Title,
			Body:        notification.Body,
			Targets:     NewNotificationTargets(notification.Targets).Response(),
			PublishedAt: notification.PublishedAt.Unix(),
			Public:      notification.Public,
			CreatedAt:   notification.CreatedAt.Unix(),
			UpdatedAt:   notification.UpdatedAt.Unix(),
		},
	}
}

func (n *Notification) Fill(administrator *Administrator) {
	if administrator != nil {
		n.CreatedBy = administrator.ID
		n.CreatorName = strings.TrimSpace(strings.Join([]string{administrator.Lastname, administrator.Firstname}, " "))
		n.UpdatedBy = administrator.ID
	}
}

func (n *Notification) Response() *response.Notification {
	return &n.Notification
}

func NewNotifications(notifications entity.Notifications) Notifications {
	res := make(Notifications, len(notifications))
	for i := range notifications {
		res[i] = NewNotification(notifications[i])
	}
	return res
}

func (ns Notifications) Response() []*response.Notification {
	res := make([]*response.Notification, len(ns))
	for i := range ns {
		res[i] = ns[i].Response()
	}
	return res
}

func (ns Notifications) AdministratorIDs() []string {
	return set.UniqBy(ns, func(n *Notification) string {
		return n.CreatedBy
	})
}

func (ns Notifications) Fill(administrators map[string]*Administrator) {
	for i := range ns {
		ns[i].Fill(administrators[ns[i].CreatedBy])
	}
}

func NewNotificationTarget(target *entity.TargetType) *TargetType {
	return &TargetType{
		TargetType: response.TargetType(*target),
	}
}

func (t *TargetType) Response() *response.TargetType {
	return &t.TargetType
}

func NewNotificationTargets(targets []entity.TargetType) TargetTypes {
	res := make(TargetTypes, len(targets))
	for i := range targets {
		res[i] = NewNotificationTarget(&targets[i])
	}
	return res
}

func (t TargetTypes) Response() []response.TargetType {
	res := make([]response.TargetType, len(t))
	for i := range t {
		res[i] = *t[i].Response()
	}
	return res
}
