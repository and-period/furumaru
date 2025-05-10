package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/validator"
	govalidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Params struct {
	WaitGroup   *sync.WaitGroup
	AdminWebURL *url.URL
	UserWebURL  *url.URL
	Database    *database.Database
	Producer    sqs.Producer
	User        user.Service
	Store       store.Service
}

type service struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	validator   validator.Validator
	adminWebURL func() *url.URL
	userWebURL  func() *url.URL
	db          *database.Database
	producer    sqs.Producer
	user        user.Service
	store       store.Service
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
	adminWebURL := func() *url.URL {
		url := *params.AdminWebURL // copy
		return &url
	}
	userWebURL := func() *url.URL {
		url := *params.UserWebURL // copy
		return &url
	}
	return &service{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		validator:   validator.NewValidator(),
		producer:    params.Producer,
		adminWebURL: adminWebURL,
		userWebURL:  userWebURL,
		db:          params.Database,
		user:        params.User,
		store:       params.Store,
	}
}

func internalError(err error) error {
	if err == nil {
		return nil
	}

	//nolint:errorlint
	if e, ok := err.(govalidator.ValidationErrors); ok {
		return fmt.Errorf("%w: %s", exception.ErrInvalidArgument, e.Error())
	}
	if e := dbError(err); e != nil {
		return fmt.Errorf("%w: %s", e, err.Error())
	}

	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", exception.ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", exception.ErrDeadlineExceeded, err.Error())
	default:
		return fmt.Errorf("%w: %s", exception.ErrInternal, err.Error())
	}
}

func dbError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, database.ErrNotFound):
		return exception.ErrNotFound
	case errors.Is(err, database.ErrFailedPrecondition):
		return exception.ErrFailedPrecondition
	case errors.Is(err, database.ErrAlreadyExists):
		return exception.ErrAlreadyExists
	case errors.Is(err, database.ErrDeadlineExceeded):
		return exception.ErrDeadlineExceeded
	default:
		return nil
	}
}
