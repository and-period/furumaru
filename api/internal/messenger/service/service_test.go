package service

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"testing"
	"time"

	mock_mailer "github.com/and-period/furumaru/api/mock/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	errmock        = errors.New("some error")
	adminWebURL, _ = url.Parse("https://admin.and-period.jp")
	userWebURL, _  = url.Parse("https://user.and-period.jp")
)

type mocks struct {
	mailer *mock_mailer.MockClient
}

type testOptions struct {
	now func() time.Time
}

type testOption func(opts *testOptions)

func withNow(now time.Time) testOption {
	return func(opts *testOptions) {
		opts.now = func() time.Time {
			return now
		}
	}
}

type testCaller func(ctx context.Context, t *testing.T, service *messengerService)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		mailer: mock_mailer.NewMockClient(ctrl),
	}
}

func newMessengerService(mocks *mocks, opts ...testOption) *messengerService {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &messengerService{
		now:       dopts.now,
		logger:    zap.NewNop(),
		waitGroup: &sync.WaitGroup{},
		validator: validator.NewValidator(),
		mailer:    mocks.mailer,
		adminWebURL: func() *url.URL {
			url := *adminWebURL
			return &url
		},
		userWebURL: func() *url.URL {
			url := *userWebURL
			return &url
		},
		maxRetries: 3,
	}
}

func testService(
	setup func(ctx context.Context, mocks *mocks),
	testFunc testCaller,
	opts ...testOption,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mocks := newMocks(ctrl)

		srv := newMessengerService(mocks)
		setup(ctx, mocks)

		testFunc(ctx, t, srv)
		srv.waitGroup.Wait()
	}
}

func TestMessengerService(t *testing.T) {
	t.Parallel()
	srv := NewMessengerService(&Params{}, WithLogger(zap.NewNop()), WithMaxRetries(3))
	assert.NotNil(t, srv)
}
