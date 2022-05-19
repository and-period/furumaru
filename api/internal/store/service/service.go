package service

import (
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/validator"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type Params struct {
	Database  *database.Database
	WaitGroup *sync.WaitGroup
}

type storeService struct {
	now         func() time.Time
	logger      *zap.Logger
	sharedGroup *singleflight.Group
	waitGroup   *sync.WaitGroup
	validator   validator.Validator
	db          *database.Database
}

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewStoreService(params *Params, opts ...Option) store.StoreService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &storeService{
		now:         jst.Now,
		logger:      dopts.logger,
		sharedGroup: &singleflight.Group{},
		validator:   validator.NewValidator(),
		waitGroup:   params.WaitGroup,
		db:          params.Database,
	}
}
