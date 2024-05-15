package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/sqs"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/and-period/furumaru/api/pkg/validator"
	"github.com/and-period/furumaru/api/pkg/youtube"
	govalidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	defaultUploadEventTTL = 12 * time.Hour     // 12hours
	defaultAuthYoutubeTTL = 3 * 24 * time.Hour // 3days
)

type Params struct {
	WaitGroup *sync.WaitGroup
	Database  *database.Database
	Cache     dynamodb.Client
	MediaLive medialive.MediaLive
	Tmp       storage.Bucket
	Storage   storage.Bucket
	Producer  sqs.Producer
	Store     store.Service
	YouTube   youtube.YouTube
}

type service struct {
	logger         *zap.Logger
	waitGroup      *sync.WaitGroup
	validator      validator.Validator
	db             *database.Database
	cache          dynamodb.Client
	tmp            storage.Bucket
	storage        storage.Bucket
	tmpURL         func() *url.URL
	storageURL     func() *url.URL
	producer       sqs.Producer
	store          store.Service
	media          medialive.MediaLive
	youtube        youtube.YouTube
	now            func() time.Time
	generateID     func() string
	uploadEventTTL time.Duration
	authYoutubeTTL time.Duration
}

type options struct {
	logger         *zap.Logger
	uploadEventTTL time.Duration
	authYoutubeTTL time.Duration
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithUploadEventTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.uploadEventTTL = ttl
	}
}

func WithAuthYoutubeTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.authYoutubeTTL = ttl
	}
}

func NewService(params *Params, opts ...Option) (media.Service, error) {
	dopts := &options{
		logger:         zap.NewNop(),
		uploadEventTTL: defaultUploadEventTTL,
		authYoutubeTTL: defaultAuthYoutubeTTL,
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
		cache:      params.Cache,
		media:      params.MediaLive,
		tmp:        params.Tmp,
		tmpURL:     tmpURL,
		storage:    params.Storage,
		storageURL: storageURL,
		producer:   params.Producer,
		store:      params.Store,
		youtube:    params.YouTube,
		now:        jst.Now,
		generateID: func() string {
			return uuid.Base58Encode(uuid.New())
		},
		uploadEventTTL: dopts.uploadEventTTL,
		authYoutubeTTL: dopts.authYoutubeTTL,
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
	if e := youtubeError(err); e != nil {
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

func youtubeError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, youtube.ErrBadRequest):
		return exception.ErrInvalidArgument
	case errors.Is(err, youtube.ErrUnauthorized):
		return exception.ErrUnauthenticated
	case errors.Is(err, youtube.ErrForbidden):
		return exception.ErrForbidden
	case errors.Is(err, youtube.ErrTooManyRequests):
		return exception.ErrResourceExhausted
	default:
		return nil
	}
}
