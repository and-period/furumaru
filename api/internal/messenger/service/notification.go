package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
)

func (s *service) CreateNotification(ctx context.Context, in *messenger.CreateNotificationInput) (*entity.Notification, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	adminId := &user.GetAdminInput{
		AdminID: in.CreatedBy,
	}
	admin, err := s.user.GetAdmin(ctx, adminId)
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid admin id format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}

	if err != nil {
		return nil, exception.InternalError(err)
	}

	params := &entity.NewNotificationParams{
		CreatedBy:   admin.ID,
		CreatorName: admin.Name(),
		UpdatedBy:   admin.ID,
		Title:       in.Title,
		Body:        in.Body,
		Targets:     in.Targets,
		Public:      in.Public,
	}
	notification := entity.NewNotification(params)
	if err := s.db.Notification.Create(ctx, notification); err != nil {
		return nil, exception.InternalError(err)
	}
	return notification, nil
}
