package service

import (
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/validator"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type Params struct {
	WaitGroup *sync.WaitGroup
	Database  *database.Database
	AdminAuth cognito.Client
	UserAuth  cognito.Client
	Messenger messenger.Service
}

type service struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	sharedGroup *singleflight.Group
	validator   validator.Validator
	db          *database.Database
	adminAuth   cognito.Client
	userAuth    cognito.Client
	messenger   messenger.Service
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

func NewService(params *Params, opts ...Option) user.Service {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &service{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		sharedGroup: &singleflight.Group{},
		validator:   validator.NewValidator(),
		db:          params.Database,
		adminAuth:   params.AdminAuth,
		userAuth:    params.UserAuth,
		messenger:   params.Messenger,
	}
}
