package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Notification struct {
	response.Notification
}

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

func (n *Notification) Response() *response.Notification {
	return &n.Notification
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
