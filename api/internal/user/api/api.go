package api

import (
	"errors"
	"sync"
	"time"

	"github.com/and-period/marche/api/internal/user/database"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/and-period/marche/api/proto/user"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Params struct {
	Logger    *zap.Logger
	WaitGroup *sync.WaitGroup
	Database  *database.Database
	UserAuth  cognito.Client
}

type userService struct {
	user.UnimplementedUserServiceServer
	now         func() time.Time
	logger      *zap.Logger
	sharedGroup *singleflight.Group
	waitGroup   *sync.WaitGroup
	db          *database.Database
	userAuth    cognito.Client
}

func NewUserService(params *Params) user.UserServiceServer {
	return &userService{
		now:         jst.Now,
		logger:      params.Logger,
		sharedGroup: &singleflight.Group{},
		waitGroup:   params.WaitGroup,
		db:          params.Database,
		userAuth:    params.UserAuth,
	}
}

func gRPCError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	switch {
	case errors.Is(err, database.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, database.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, database.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, database.ErrNotImplemented):
		return status.Error(codes.Unimplemented, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
