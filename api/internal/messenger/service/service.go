package service

import (
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/validator"
	"go.uber.org/zap"
)

type Params struct {
	WaitGroup *sync.WaitGroup
	Producer  sqs.Producer
}

type service struct {
	now       func() time.Time
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
	validator validator.Validator
	producer  sqs.Producer
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

func NewService(params *Params, opts ...Option) messenger.Service {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &service{
		now:       jst.Now,
		logger:    dopts.logger,
		waitGroup: params.WaitGroup,
		validator: validator.NewValidator(),
		producer:  params.Producer,
	}
}
