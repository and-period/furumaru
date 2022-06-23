package service

import (
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/validator"
	"go.uber.org/zap"
)

type Params struct {
	WaitGroup   *sync.WaitGroup
	Mailer      mailer.Client
	AdminWebURL *url.URL
	UserWebURL  *url.URL
}

type service struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	validator   validator.Validator
	mailer      mailer.Client
	adminWebURL func() *url.URL
	userWebURL  func() *url.URL
	maxRetries  int64
}

type options struct {
	logger     *zap.Logger
	maxRetries int64
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithMaxRetries(maxRetries int64) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func NewService(params *Params, opts ...Option) messenger.Service {
	dopts := &options{
		logger:     zap.NewNop(),
		maxRetries: 3,
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
		mailer:      params.Mailer,
		adminWebURL: adminWebURL,
		userWebURL:  userWebURL,
		maxRetries:  dopts.maxRetries,
	}
}
