package api

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/and-period/marche/api/internal/user/database"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var errmock = errors.New("some error")

type mocks struct {
	db *dbMocks
}

type dbMocks struct{}

type testResponse struct {
	code codes.Code
	body proto.Message
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

type grpcCaller func(ctx context.Context, service *userService) (proto.Message, error)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db: newDBMocks(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{}
}

func newUserService(mocks *mocks, opts *testOptions) *userService {
	return &userService{
		now:         jst.Now,
		logger:      zap.NewNop(),
		sharedGroup: &singleflight.Group{},
		waitGroup:   &sync.WaitGroup{},
		db:          &database.Database{},
	}
}

func testGRPC(
	setup func(context.Context, *testing.T, *mocks),
	expect *testResponse,
	grpcFn grpcCaller,
	opts ...testOption,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dopts := &testOptions{
			now: jst.Now,
		}
		for i := range opts {
			opts[i](dopts)
		}

		mocks := newMocks(ctrl)
		service := newUserService(mocks, dopts)
		if setup != nil {
			setup(ctx, t, mocks)
		}

		res, err := grpcFn(ctx, service)
		if expect != nil {
			switch expect.code {
			case codes.OK:
				require.NoError(t, err)
			default:
				require.Error(t, err)
				status, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, expect.code, status.Code(), status.Code().String())
			}
			if expect.body != nil {
				require.Equal(t, expect.body, res)
			}
		}
		service.waitGroup.Wait()
	}
}

func TestUserService(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, NewUserService(&Params{}))
}

func TestGRPCError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect codes.Code
	}{
		{
			name:   "error is nil",
			err:    nil,
			expect: codes.OK,
		},
		{
			name:   "grpc error",
			err:    status.Error(codes.Unavailable, errmock.Error()),
			expect: codes.Unavailable,
		},
		{
			name:   "invalid argument",
			err:    fmt.Errorf("%w: %s", database.ErrInvalidArgument, errmock),
			expect: codes.InvalidArgument,
		},
		{
			name:   "not found",
			err:    fmt.Errorf("%w: %s", database.ErrNotFound, errmock),
			expect: codes.NotFound,
		},
		{
			name:   "already exists",
			err:    fmt.Errorf("%w: %s", database.ErrAlreadyExists, errmock),
			expect: codes.AlreadyExists,
		},
		{
			name:   "unimplemented",
			err:    fmt.Errorf("%w: %s", database.ErrNotImplemented, errmock),
			expect: codes.Unimplemented,
		},
		{
			name:   "internal",
			err:    fmt.Errorf("%w: %s", database.ErrInternal, errmock),
			expect: codes.Internal,
		},
		{
			name:   "other error",
			err:    errmock,
			expect: codes.Unknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := gRPCError(tt.err)
			assert.Equal(t, tt.expect, status.Code(err))
		})
	}
}
