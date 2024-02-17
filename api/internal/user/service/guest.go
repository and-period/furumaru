package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (s *service) UpsertGuest(ctx context.Context, in *user.UpsertGuestInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.User{}, nil
}
