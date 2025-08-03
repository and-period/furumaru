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
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/codes"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/validator"
	govalidator "github.com/go-playground/validator/v10"
	"golang.org/x/sync/singleflight"
)

const (
	defaultAdminAuthTTL = 5 * time.Minute
	defaultUserAuthTTL  = 5 * time.Minute
)

type Params struct {
	WaitGroup                  *sync.WaitGroup
	Database                   *database.Database
	Cache                      dynamodb.Client
	AdminAuth                  cognito.Client
	UserAuth                   cognito.Client
	Store                      store.Service
	Messenger                  messenger.Service
	Media                      media.Service
	DefaultAdminGroups         map[entity.AdminType][]string
	AdminAuthGoogleRedirectURL string
	AdminAuthLINERedirectURL   string
	UserAuthGoogleRedirectURL  string
	UserAuthLINERedirectURL    string
}

type service struct {
	now                        func() time.Time
	waitGroup                  *sync.WaitGroup
	sharedGroup                *singleflight.Group
	validator                  validator.Validator
	db                         *database.Database
	cache                      dynamodb.Client
	adminAuth                  cognito.Client
	userAuth                   cognito.Client
	store                      store.Service
	messenger                  messenger.Service
	media                      media.Service
	defaultAdminGroups         map[entity.AdminType][]string
	adminAuthTTL               time.Duration
	adminAuthGoogleRedirectURL string
	adminAuthLINERedirectURL   string
	userAuthTTL                time.Duration
	userAuthGoogleRedirectURL  string
	userAuthLINERedirectURL    string
}

type options struct {
	adminAuthTTL time.Duration
	userAuthTTL  time.Duration
}

type Option func(*options)

func WithAdminAuthTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.adminAuthTTL = ttl
	}
}

func WithUserAuthTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.userAuthTTL = ttl
	}
}

func NewService(params *Params, opts ...Option) user.Service {
	dopts := &options{
		adminAuthTTL: defaultAdminAuthTTL,
		userAuthTTL:  defaultUserAuthTTL,
	}
	for i := range opts {
		opts[i](dopts)
	}
	vopts := []validator.Option{
		validator.WithPasswordValidation(&validator.PasswordParams{
			RequireNumbers:   true,
			RequireSymbols:   false,
			RequireUppercase: false,
			RequireLowercase: true,
		}),
		validator.WithCustomValidation(codes.RegisterValidations),
	}
	return &service{
		now:                        jst.Now,
		waitGroup:                  params.WaitGroup,
		sharedGroup:                &singleflight.Group{},
		validator:                  validator.NewValidator(vopts...),
		db:                         params.Database,
		cache:                      params.Cache,
		adminAuth:                  params.AdminAuth,
		userAuth:                   params.UserAuth,
		store:                      params.Store,
		messenger:                  params.Messenger,
		media:                      params.Media,
		defaultAdminGroups:         params.DefaultAdminGroups,
		adminAuthTTL:               dopts.adminAuthTTL,
		adminAuthGoogleRedirectURL: params.AdminAuthGoogleRedirectURL,
		adminAuthLINERedirectURL:   params.AdminAuthLINERedirectURL,
		userAuthTTL:                dopts.userAuthTTL,
		userAuthGoogleRedirectURL:  params.UserAuthGoogleRedirectURL,
		userAuthLINERedirectURL:    params.UserAuthLINERedirectURL,
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
	if e := authError(err); e != nil {
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

func authError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, cognito.ErrInvalidArgument):
		return exception.ErrInvalidArgument
	case errors.Is(err, cognito.ErrUnauthenticated), errors.Is(err, cognito.ErrCodeExpired):
		return exception.ErrUnauthenticated
	case errors.Is(err, cognito.ErrNotFound):
		return exception.ErrNotFound
	case errors.Is(err, cognito.ErrAlreadyExists):
		return exception.ErrAlreadyExists
	case errors.Is(err, cognito.ErrResourceExhausted):
		return exception.ErrResourceExhausted
	case errors.Is(err, cognito.ErrTimeout):
		return exception.ErrDeadlineExceeded
	default:
		return nil
	}
}
