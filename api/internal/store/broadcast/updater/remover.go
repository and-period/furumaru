package updater

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"go.uber.org/zap"
)

type remover struct {
	now        func() time.Time
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	db         *database.Database
	maxRetries int64
}

func NewRemover(params *Params, opts ...Option) Updater {
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

func (r *remover) Lambda(ctx context.Context, event interface{}) error {
	payload, ok := event.(RemovePayload)
	if !ok {
		r.logger.Error("Received unexpected event format", zap.Any("event", event))
		return errors.New("updater: received unexpected event format")
	}
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
