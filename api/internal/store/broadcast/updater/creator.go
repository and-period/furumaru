package updater

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"go.uber.org/zap"
)

type creator struct {
	now        func() time.Time
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	db         *database.Database
	maxRetries int64
}

func NewCreator(params *Params, opts ...Option) Updater {
	dopts := &options{
		logger:     zap.NewNop(),
		maxRetries: 3,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &creator{
		now:        jst.Now,
		logger:     dopts.logger,
		waitGroup:  params.WaitGroup,
		db:         params.Database,
		maxRetries: dopts.maxRetries,
	}
}

func (c *creator) Lambda(_ context.Context, event interface{}) error {
	c.logger.Debug("Received event", zap.Any("event", event))
	// TODO: 取得内容が分かり次第、詳細の実装
	return nil
}
