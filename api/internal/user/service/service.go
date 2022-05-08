//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/user/$GOPACKAGE/$GOFILE
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/marche/api/internal/user/database"
	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/and-period/marche/api/pkg/storage"
	"github.com/and-period/marche/api/pkg/validator"
	gvalidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var (
	ErrInvalidArgument    = errors.New("service: invalid argument")
	ErrUnauthenticated    = errors.New("service: unauthenticated")
	ErrNotFound           = errors.New("service: not found")
	ErrAlreadyExists      = errors.New("service: already exists")
	ErrFailedPrecondition = errors.New("service: failed precondition")
	ErrResourceExhausted  = errors.New("service: resource exhausted")
	ErrNotImplemented     = errors.New("service: not implemented")
	ErrInternal           = errors.New("service: internal error")
)

type UserService interface {
	SignInAdmin(ctx context.Context, in *SignInAdminInput) (*entity.AdminAuth, error)
	SignOutAdmin(ctx context.Context, in *SignOutAdminInput) error
	GetAdminAuth(ctx context.Context, in *GetAdminAuthInput) (*entity.AdminAuth, error)
	RefreshAdminToken(ctx context.Context, in *RefreshAdminTokenInput) (*entity.AdminAuth, error)
	GetAdmin(ctx context.Context, in *GetAdminInput) (*entity.Admin, error)
	UpdateAdminEmail(ctx context.Context, in *UpdateAdminEmailInput) error
	VerifyAdminEmail(ctx context.Context, in *VerifyAdminEmailInput) error
	UpdateAdminPassword(ctx context.Context, in *UpdateAdminPasswordInput) error
	MultiGetShops(ctx context.Context, in *MultiGetShopsInput) (entity.Shops, error)
	GetShop(ctx context.Context, in *GetShopInput) (*entity.Shop, error)
	SignInUser(ctx context.Context, in *SignInUserInput) (*entity.UserAuth, error)
	SignOutUser(ctx context.Context, in *SignOutUserInput) error
	GetUserAuth(ctx context.Context, in *GetUserAuthInput) (*entity.UserAuth, error)
	RefreshUserToken(ctx context.Context, in *RefreshUserTokenInput) (*entity.UserAuth, error)
	GetUser(ctx context.Context, in *GetUserInput) (*entity.User, error)
	CreateUser(ctx context.Context, in *CreateUserInput) (string, error)
	VerifyUser(ctx context.Context, in *VerifyUserInput) error
	CreateUserWithOAuth(ctx context.Context, in *CreateUserWithOAuthInput) (*entity.User, error)
	InitializeUser(ctx context.Context, in *InitializeUserInput) error
	UpdateUserEmail(ctx context.Context, in *UpdateUserEmailInput) error
	VerifyUserEmail(ctx context.Context, in *VerifyUserEmailInput) error
	UpdateUserPassword(ctx context.Context, in *UpdateUserPasswordInput) error
	ForgotUserPassword(ctx context.Context, in *ForgotUserPasswordInput) error
	VerifyUserPassword(ctx context.Context, in *VerifyUserPasswordInput) error
	DeleteUser(ctx context.Context, in *DeleteUserInput) error
}

type Params struct {
	Storage   storage.Bucket
	Database  *database.Database
	AdminAuth cognito.Client
	ShopAuth  cognito.Client
	UserAuth  cognito.Client
}

type userService struct {
	now         func() time.Time
	logger      *zap.Logger
	sharedGroup *singleflight.Group
	validator   validator.Validator
	storage     storage.Bucket
	db          *database.Database
	adminAuth   cognito.Client
	shopAuth    cognito.Client
	userAuth    cognito.Client
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

func NewUserService(params *Params, opts ...Option) UserService {
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
		storage:     params.Storage,
		db:          params.Database,
		adminAuth:   params.AdminAuth,
		shopAuth:    params.ShopAuth,
		userAuth:    params.UserAuth,
	}
}

func userError(err error) error {
	if err == nil {
		return nil
	}

	//nolint:gocritic
	switch v := err.(type) {
	case gvalidator.ValidationErrors:
		return fmt.Errorf("%w: %s", ErrInvalidArgument, v.Error())
	}

	switch {
	case errors.Is(err, database.ErrInvalidArgument), errors.Is(err, cognito.ErrInvalidArgument):
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case errors.Is(err, cognito.ErrUnauthenticated):
		return fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	case errors.Is(err, database.ErrNotFound), errors.Is(err, cognito.ErrNotFound):
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case errors.Is(err, database.ErrAlreadyExists), errors.Is(err, cognito.ErrUnauthenticated):
		return fmt.Errorf("%w: %s", ErrAlreadyExists, err.Error())
	case errors.Is(err, cognito.ErrResourceExhausted):
		return fmt.Errorf("%w: %s", ErrResourceExhausted, err.Error())
	case errors.Is(err, database.ErrNotImplemented):
		return fmt.Errorf("%w: %s", ErrNotImplemented, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	}
}
