package service

import (
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/validator"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type Params struct {
	WaitGroup  *sync.WaitGroup
	Database   *database.Database
	User       user.Service
	Messenger  messenger.Service
	Media      media.Service
	PostalCode postalcode.Client
}

type service struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	sharedGroup *singleflight.Group
	validator   validator.Validator
	db          *database.Database
	user        user.Service
	messenger   messenger.Service
	media       media.Service
	postalCode  postalcode.Client
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

func NewService(params *Params, opts ...Option) store.Service {
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
		user:        params.User,
		messenger:   params.Messenger,
		media:       params.Media,
		postalCode:  params.PostalCode,
	}
}
