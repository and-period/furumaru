package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (s *service) GetUserNotification(ctx context.Context, in *user.GetUserNotificationInput) (*entity.UserNotification, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	notification, err := s.db.UserNotification.Get(ctx, in.UserID)
	return notification, internalError(err)
}

func (s *service) UpdateUserNotification(ctx context.Context, in *user.UpdateUserNotificationInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	notification, err := s.db.UserNotification.Get(ctx, in.UserID)
	if err != nil {
		return internalError(err)
	}
	notification.EmailDisabled = !in.Enabled
	err = s.db.UserNotification.Upsert(ctx, notification)
	return internalError(err)
}
