package service

import (
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/and-period/furumaru/api/pkg/validator"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type Params struct {
	Storage          storage.Bucket
	Database         *database.Database
	AdminAuth        cognito.Client
	UserAuth         cognito.Client
	MessengerService messenger.MessengerService
	WaitGroup        *sync.WaitGroup
}

type userService struct {
	now         func() time.Time
	logger      *zap.Logger
	sharedGroup *singleflight.Group
	waitGroup   *sync.WaitGroup
	validator   validator.Validator
	storage     storage.Bucket
	db          *database.Database
	adminAuth   cognito.Client
	userAuth    cognito.Client
	messenger   messenger.MessengerService
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

func NewUserService(params *Params, opts ...Option) user.UserService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &userService{
		now:         jst.Now,
		logger:      dopts.logger,
		sharedGroup: &singleflight.Group{},
		validator:   validator.NewValidator(),
		waitGroup:   params.WaitGroup,
		storage:     params.Storage,
		db:          params.Database,
		adminAuth:   params.AdminAuth,
		userAuth:    params.UserAuth,
		messenger:   params.MessengerService,
	}
}
