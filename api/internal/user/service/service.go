package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/validator"
	govalidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type Params struct {
	WaitGroup *sync.WaitGroup
	Database  *database.Database
	AdminAuth cognito.Client
	UserAuth  cognito.Client
	Store     store.Service
	Messenger messenger.Service
	Media     media.Service
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
	store       store.Service
	messenger   messenger.Service
	media       media.Service
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
		store:       params.Store,
		messenger:   params.Messenger,
		media:       params.Media,
	}
}

func internalError(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(govalidator.ValidationErrors); ok {
		return fmt.Errorf("%w: %s", user.ErrInvalidArgument, e.Error())
	}
	if e := dbError(err); e != nil {
		return fmt.Errorf("%w: %s", e, err.Error())
	}
	if e := authError(err); e != nil {
		return fmt.Errorf("%w: %s", e, err.Error())
	}

	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", user.ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", user.ErrDeadlineExceeded, err.Error())
	default:
		return fmt.Errorf("%w: %s", user.ErrInternal, err.Error())
	}
}

func dbError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, database.ErrNotFound):
		return user.ErrNotFound
	case errors.Is(err, database.ErrFailedPrecondition):
		return user.ErrFailedPrecondition
	case errors.Is(err, database.ErrAlreadyExists):
		return user.ErrAlreadyExists
	case errors.Is(err, database.ErrDeadlineExceeded):
		return user.ErrDeadlineExceeded
	default:
		return nil
	}
}

func authError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, cognito.ErrInvalidArgument):
		return user.ErrInvalidArgument
	case errors.Is(err, cognito.ErrUnauthenticated):
		return user.ErrUnauthenticated
	case errors.Is(err, cognito.ErrNotFound):
		return user.ErrNotFound
	case errors.Is(err, cognito.ErrAlreadyExists):
		return user.ErrAlreadyExists
	case errors.Is(err, cognito.ErrResourceExhausted):
		return user.ErrResourceExhausted
	case errors.Is(err, cognito.ErrTimeout):
		return user.ErrDeadlineExceeded
	default:
		return nil
	}
}
