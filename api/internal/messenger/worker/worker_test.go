package worker

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"testing"
	"time"

	mock_mailer "github.com/and-period/furumaru/api/mock/pkg/mailer"
	mock_user "github.com/and-period/furumaru/api/mock/user"
	"github.com/and-period/furumaru/api/pkg/jst"
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
	user   *mock_user.MockService
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

type testCaller func(ctx context.Context, t *testing.T, worker *worker)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		mailer: mock_mailer.NewMockClient(ctrl),
		user:   mock_user.NewMockService(ctrl),
	}
}

func newWorker(mocks *mocks, opts ...testOption) *worker {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	adminWebURL := func() *url.URL {
		url := *adminWebURL // copy
		return &url
	}
	userWebURL := func() *url.URL {
		url := *userWebURL // copy
		return &url
	}
	return &worker{
		now:         dopts.now,
		logger:      zap.NewNop(),
		waitGroup:   &sync.WaitGroup{},
		mailer:      mocks.mailer,
		adminWebURL: adminWebURL,
		userWebURL:  userWebURL,
		user:        mocks.user,
		concurrency: 1,
		maxRetries:  1,
	}
}

func testWorker(
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

		w := newWorker(mocks)
		setup(ctx, mocks)

		testFunc(ctx, t, w)
		w.waitGroup.Wait()
	}
}

func TestWorker(t *testing.T) {
	t.Parallel()
	w := NewWorker(&Params{}, WithLogger(zap.NewNop()), WithConcurrency(1), WithMaxRetries(3))
	assert.NotNil(t, w)
}
