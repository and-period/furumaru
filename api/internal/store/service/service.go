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
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/payment"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/geolocation"
	"github.com/and-period/furumaru/api/pkg/ivs"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/and-period/furumaru/api/pkg/validator"
	govalidator "github.com/go-playground/validator/v10"
	"golang.org/x/sync/singleflight"
)

var errUnmatchProducts = errors.New("service: umnatch products")

const (
	defaultCartTTL             = 14 * 24 * time.Hour // 14days
	defaultCartRefreshInterval = 2 * time.Hour       // 2hours
)

type Params struct {
	WaitGroup   *sync.WaitGroup
	Database    *database.Database
	Cache       dynamodb.Client
	User        user.Service
	Messenger   messenger.Service
	Media       media.Service
	PostalCode  postalcode.Client
	Geolocation geolocation.Client
	Ivs         ivs.Client
	Providers   map[entity.PaymentProviderType]payment.Provider
}

type service struct {
	now                 func() time.Time
	generateID          func() string
	waitGroup           *sync.WaitGroup
	sharedGroup         *singleflight.Group
	validator           validator.Validator
	db                  *database.Database
	cache               dynamodb.Client
	user                user.Service
	messenger           messenger.Service
	media               media.Service
	postalCode          postalcode.Client
	geolocation         geolocation.Client
	ivs                 ivs.Client
	providers           map[entity.PaymentProviderType]payment.Provider
	cartTTL             time.Duration
	cartRefreshInterval time.Duration
}

type options struct {
	cartTTL             time.Duration
	cartRefreshInterval time.Duration
}

type Option func(*options)

func WithCartTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.cartTTL = ttl
	}
}

func WithCartRefreshInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.cartRefreshInterval = interval
	}
}

func NewService(params *Params, opts ...Option) store.Service {
	dopts := &options{
		cartTTL:             defaultCartTTL,
		cartRefreshInterval: defaultCartRefreshInterval,
	}
	for i := range opts {
		opts[i](dopts)
	}
	providers := params.Providers
	if providers == nil {
		providers = make(map[entity.PaymentProviderType]payment.Provider)
	}
	return &service{
		now: jst.Now,
		generateID: func() string {
			return uuid.Base58Encode(uuid.New())
		},
		waitGroup:           params.WaitGroup,
		sharedGroup:         &singleflight.Group{},
		validator:           validator.NewValidator(),
		db:                  params.Database,
		cache:               params.Cache,
		user:                params.User,
		messenger:           params.Messenger,
		media:               params.Media,
		postalCode:          params.PostalCode,
		geolocation:         params.Geolocation,
		ivs:                 params.Ivs,
		providers:           providers,
		cartTTL:             dopts.cartTTL,
		cartRefreshInterval: defaultCartRefreshInterval,
	}
}

// getProvider は決済手段に対応するプロバイダーを返す。
// PaymentSystem テーブルの ProviderType に基づいて適切なプロバイダーを選択する。
func (s *service) getProvider(ctx context.Context, methodType entity.PaymentMethodType) (payment.Provider, entity.PaymentProviderType, error) {
	system, err := s.db.PaymentSystem.Get(ctx, methodType)
	if err != nil {
		return nil, entity.PaymentProviderTypeUnknown, err
	}
	prov, ok := s.providers[system.ProviderType]
	if !ok {
		return nil, entity.PaymentProviderTypeUnknown, fmt.Errorf("service: unsupported provider type: %d: %w", system.ProviderType, exception.ErrInternal)
	}
	return prov, system.ProviderType, nil
}

// getProviderByType はプロバイダー種別から直接プロバイダーを返す。
// 既に注文に記録されたプロバイダー種別から逆引きする際に使用する。
func (s *service) getProviderByType(providerType entity.PaymentProviderType) (payment.Provider, error) {
	prov, ok := s.providers[providerType]
	if !ok {
		return nil, fmt.Errorf("service: unsupported provider type: %d: %w", providerType, exception.ErrInternal)
	}
	return prov, nil
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

	//nolint:errorlint
	if e, ok := err.(govalidator.ValidationErrors); ok {
		return fmt.Errorf("%w: %s", exception.ErrInvalidArgument, e.Error())
	}
	if e := dbError(err); e != nil {
		return fmt.Errorf("%w: %s", e, err.Error())
	}
	if e := cacheError(err); e != nil {
		return fmt.Errorf("%w: %s", e, err.Error())
	}
	if e := postalCodeError(err); e != nil {
		return fmt.Errorf("%w: %s", e, err.Error())
	}
	if e := geolocationError(err); e != nil {
		return fmt.Errorf("%w: %s", e, err.Error())
	}
	if e := entityError(err); e != nil {
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

func cacheError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, dynamodb.ErrNotFound):
		return exception.ErrNotFound
	case errors.Is(err, dynamodb.ErrAlreadyExists):
		return exception.ErrAlreadyExists
	case errors.Is(err, dynamodb.ErrResourceExhausted), errors.Is(err, dynamodb.ErrOutOfRange):
		return exception.ErrResourceExhausted
	case errors.Is(err, dynamodb.ErrAborted), errors.Is(err, dynamodb.ErrCanceled):
		return exception.ErrCanceled
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

func geolocationError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, geolocation.ErrNotFound):
		return exception.ErrInvalidArgument
	default:
		return nil
	}
}

func entityError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, entity.ErrNotFoundShippingRate),
		errors.Is(err, entity.ErrUnknownShippingSize):
		return exception.ErrFailedPrecondition
	default:
		return nil
	}
}
