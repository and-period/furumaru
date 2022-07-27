package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Notification struct {
	response.Notification
}

func NewNotification(notification *entity.Notification) *Notification {
	return &Notification{
		Notification: response.Notification{
			ID:          notification.ID,
			CreatedBy:   notification.CreatedBy,
			CreatorName: notification.CreatorName,
			UpdatedBy:   notification.UpdatedBy,
			Title:       notification.Title,
			Body:        notification.Body,
			Targets:     notification.Targets,
			PublishedAt: notification.PublishedAt.Unix(),
			Public:      notification.Public,
			CreatedAt:   notification.CreatedAt.Unix(),
			UpdatedAt:   notification.UpdatedAt.Unix(),
		},
	}
}
