package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *service) GetThread(ctx context.Context, in *messenger.GetThreadInput) (*entity.Thread, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}

	thread, err := s.db.Thread.Get(ctx, in.ThreadID)
	if err != nil {
		return nil, exception.InternalError(err)
	}

	return thread, exception.InternalError(err)
}
