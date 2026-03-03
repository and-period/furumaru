package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListAuditLogs(
	ctx context.Context, in *user.ListAuditLogsInput,
) (entity.AuditLogs, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListAuditLogsParams{
		AdminID:      in.AdminID,
		ResourceType: in.ResourceType,
		Action:       in.Action,
		StartAt:      in.StartAt,
		EndAt:        in.EndAt,
		Limit:        int(in.Limit),
		Offset:       int(in.Offset),
	}

	var (
		logs  entity.AuditLogs
		total int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		logs, err = s.db.AuditLog.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.AuditLog.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return logs, total, nil
}
