package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/ivs"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/validator"
	govalidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var errUnmatchProducts = errors.New("service: umnatch products")

type Params struct {
	WaitGroup  *sync.WaitGroup
	Database   *database.Database
	User       user.Service
	Messenger  messenger.Service
	Media      media.Service
	PostalCode postalcode.Client
	Ivs        ivs.Client
	Komoju     *komoju.Komoju
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
	ivs         ivs.Client
	komoju      *komoju.Komoju
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
		ivs:         params.Ivs,
		komoju:      params.Komoju,
	}
}

func (s *service) isRetryable(err error) bool {
	return errors.Is(err, exception.ErrDeadlineExceeded) ||
		errors.Is(err, exception.ErrInternal) ||
		errors.Is(err, exception.ErrDeadlineExceeded) ||
		errors.Is(err, exception.ErrInternal)
}

func internalError(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(govalidator.ValidationErrors); ok {
		return fmt.Errorf("%w: %s", exception.ErrInvalidArgument, e.Error())
	}
	if e := dbError(err); e != nil {
		return fmt.Errorf("%w: %s", e, err.Error())
	}
	if e := postalCodeError(err); e != nil {
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

func postalCodeError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, postalcode.ErrInvalidArgument):
		return exception.ErrInvalidArgument
	case errors.Is(err, postalcode.ErrNotFound):
		return exception.ErrNotFound
	case errors.Is(err, postalcode.ErrUnavailable):
		return exception.ErrUnavailable
	case errors.Is(err, postalcode.ErrTimeout):
		return exception.ErrDeadlineExceeded
	default:
		return nil
	}
}
