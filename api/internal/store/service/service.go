//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/marche/api/internal/store/database"
	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/jst"
	validator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var (
	ErrInvalidArgument    = errors.New("service: invalid argument")
	ErrUnauthenticated    = errors.New("service: unauthenticated")
	ErrNotFound           = errors.New("service: not found")
	ErrAlreadyExists      = errors.New("service: already exists")
	ErrFailedPrecondition = errors.New("service: failed precondition")
	ErrNotImplemented     = errors.New("service: not implemented")
	ErrInternal           = errors.New("service: internal error")
)

type StoreService interface {
	ListStaffsByStoreID(ctx context.Context, in *ListStaffsByStoreIDInput) (entity.Staffs, error)
	ListStores(ctx context.Context, in *ListStoresInput) (entity.Stores, error)
	GetStore(ctx context.Context, in *GetStoreInput) (*entity.Store, error)
}

type Params struct {
	Database *database.Database
}

type storeService struct {
	now         func() time.Time
	logger      *zap.Logger
	sharedGroup *singleflight.Group
	validator   *validator.Validate
	db          *database.Database
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

func NewStoreService(params *Params, opts ...Option) StoreService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &storeService{
		now:         jst.Now,
		logger:      dopts.logger,
		sharedGroup: &singleflight.Group{},
		validator:   newValidator(),
		db:          params.Database,
	}
}

func storeError(err error) error {
	if err == nil {
		return nil
	}

	//nolint:gocritic
	switch v := err.(type) {
	case validator.ValidationErrors:
		return fmt.Errorf("%w: %s", ErrInvalidArgument, v.Error())
	}

	switch {
	case errors.Is(err, database.ErrInvalidArgument):
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case errors.Is(err, database.ErrNotFound):
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case errors.Is(err, database.ErrAlreadyExists):
		return fmt.Errorf("%w: %s", ErrAlreadyExists, err.Error())
	case errors.Is(err, database.ErrNotImplemented):
		return fmt.Errorf("%w: %s", ErrNotImplemented, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	}
}
