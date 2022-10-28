package service

import (
	"sync"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/and-period/furumaru/api/pkg/validator"
	"go.uber.org/zap"
)

type Params struct {
	WaitGroup *sync.WaitGroup
	Temp      storage.Bucket
	Storage   storage.Bucket
}

type service struct {
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	validator   validator.Validator
	temp        storage.Bucket
	storage     storage.Bucket
	tempHost    string
	storageHost string
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

func NewService(params *Params, opts ...Option) media.Service {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &service{
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		validator:   validator.NewValidator(),
		temp:        params.Temp,
		tempHost:    params.Temp.GetFQDN(),
		storage:     params.Storage,
		storageHost: params.Temp.GetFQDN(),
	}
}
