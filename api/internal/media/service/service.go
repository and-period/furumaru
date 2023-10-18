package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/and-period/furumaru/api/pkg/validator"
	govalidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	errParseURL   = errors.New("service: failed to parse url")
	errInvalidURL = errors.New("service: invalid url")
)

type Params struct {
	WaitGroup *sync.WaitGroup
	Database  *database.Database
	Tmp       storage.Bucket
	Storage   storage.Bucket
	Producer  sqs.Producer
}

type service struct {
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	validator  validator.Validator
	db         *database.Database
	tmp        storage.Bucket
	storage    storage.Bucket
	tmpURL     func() *url.URL
	storageURL func() *url.URL
	producer   sqs.Producer
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

func NewService(params *Params, opts ...Option) (media.Service, error) {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	turl, err := params.Tmp.GetHost()
	if err != nil {
		return nil, err
	}
	surl, err := params.Storage.GetHost()
	if err != nil {
		return nil, err
	}
	tmpURL := func() *url.URL {
		url := *turl // copy
		return &url
	}
	storageURL := func() *url.URL {
		url := *surl // copy
		return &url
	}
	return &service{
		logger:     dopts.logger,
		waitGroup:  params.WaitGroup,
		validator:  validator.NewValidator(),
		db:         params.Database,
		tmp:        params.Tmp,
		tmpURL:     tmpURL,
		storage:    params.Storage,
		storageURL: storageURL,
		producer:   params.Producer,
	}, nil
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
	if e := storageError(err); e != nil {
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

func storageError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, storage.ErrInvalidURL):
		return exception.ErrInvalidArgument
	case errors.Is(err, storage.ErrNotFound):
		return exception.ErrNotFound
	default:
		return nil
	}
}
