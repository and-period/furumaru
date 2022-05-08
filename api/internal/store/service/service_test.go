package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/and-period/marche/api/internal/store/database"
	mock_database "github.com/and-period/marche/api/mock/store/database"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/and-period/marche/api/pkg/validator"
	gvalidator "github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var errmock = errors.New("some error")

type mocks struct {
	db *dbMocks
}

type dbMocks struct {
	Staff *mock_database.MockStaff
	Store *mock_database.MockStore
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

type testCaller func(ctx context.Context, t *testing.T, service *storeService)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db: newDBMocks(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Staff: mock_database.NewMockStaff(ctrl),
		Store: mock_database.NewMockStore(ctrl),
	}
}

func newUserService(mocks *mocks, opts ...testOption) *storeService {
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &storeService{
		now:         dopts.now,
		logger:      zap.NewNop(),
		sharedGroup: &singleflight.Group{},
		validator:   validator.NewValidator(),
		db: &database.Database{
			Staff: mocks.db.Staff,
			Store: mocks.db.Store,
		},
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

		srv := newUserService(mocks)
		setup(ctx, mocks)

		testFunc(ctx, t, srv)
	}
}

func TestStoreService(t *testing.T) {
	t.Parallel()
	srv := NewStoreService(&Params{}, WithLogger(zap.NewNop()))
	assert.NotNil(t, srv)
}

func TestStoreError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect error
	}{
		{
			name:   "error is nil",
			err:    nil,
			expect: nil,
		},
		{
			name:   "validation error",
			err:    gvalidator.ValidationErrors{},
			expect: ErrInvalidArgument,
		},
		{
			name:   "invalid argument",
			err:    fmt.Errorf("%w: %s", database.ErrInvalidArgument, errmock),
			expect: ErrInvalidArgument,
		},
		{
			name:   "not found",
			err:    fmt.Errorf("%w: %s", database.ErrNotFound, errmock),
			expect: ErrNotFound,
		},
		{
			name:   "already exists",
			err:    fmt.Errorf("%w: %s", database.ErrAlreadyExists, errmock),
			expect: ErrAlreadyExists,
		},
		{
			name:   "unimplemented",
			err:    fmt.Errorf("%w: %s", database.ErrNotImplemented, errmock),
			expect: ErrNotImplemented,
		},
		{
			name:   "other error",
			err:    errmock,
			expect: ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := storeError(tt.err)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}
