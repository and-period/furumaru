package updater

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"go.uber.org/zap"
)

type Remover interface {
	Lambda(ctx context.Context, event RemovePayload) error
}

type remover struct {
	now        func() time.Time
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	db         *database.Database
	maxRetries int64
}

func NewRemover(params *Params, opts ...Option) Remover {
	dopts := &options{
		logger:     zap.NewNop(),
		maxRetries: 3,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &remover{
		now:        jst.Now,
		logger:     dopts.logger,
		waitGroup:  params.WaitGroup,
		db:         params.Database,
		maxRetries: dopts.maxRetries,
	}
}

func (r *remover) Lambda(ctx context.Context, payload RemovePayload) error {
	r.logger.Debug("Received event", zap.Any("event", payload))
	params := &database.UpdateBroadcastParams{
		Status: entity.BroadcastStatusDisabled,
	}
	if err := r.db.Broadcast.Update(ctx, payload.ScheduleID, params); err != nil {
		r.logger.Error("Failed to update broadcast", zap.Error(err), zap.String("scheduleId", payload.ScheduleID))
		return err
	}
	r.logger.Info("Succeeded to update broadcast", zap.String("scheduleId", payload.ScheduleID))
	return nil
}
