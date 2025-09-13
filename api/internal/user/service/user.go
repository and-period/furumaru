package service

import (
	"context"
	"errors"

	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListUsers(ctx context.Context, in *user.ListUsersInput) (entity.Users, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListUsersParams{
		Limit:          int(in.Limit),
		Offset:         int(in.Offset),
		OnlyRegistered: in.OnlyRegistered,
		OnlyVerified:   in.OnlyVerified,
		WithDeleted:    in.WithDeleted,
	}
	var (
		users entity.Users
		total int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		users, err = s.db.User.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.User.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return users, total, nil
}

func (s *service) MultiGetUsers(ctx context.Context, in *user.MultiGetUsersInput) (entity.Users, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	users, err := s.db.User.MultiGet(ctx, in.UserIDs)
	return users, internalError(err)
}

func (s *service) MultiGetUserDevices(_ context.Context, in *user.MultiGetUserDevicesInput) ([]string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return []string{}, nil
}

func (s *service) GetUser(ctx context.Context, in *user.GetUserInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	return u, internalError(err)
}

func (s *service) DeleteUser(ctx context.Context, in *user.DeleteUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	if err != nil {
		return internalError(err)
	}
	switch u.Type {
	case entity.UserTypeMember:
		auth := func(ctx context.Context) error {
			err := s.userAuth.DeleteUser(ctx, u.CognitoID)
			if errors.Is(err, cognito.ErrNotFound) {
				return nil // すでに削除済み
			}
			return err
		}
		err = s.db.Member.Delete(ctx, u.ID, auth)
	case entity.UserTypeGuest:
		err = s.db.Guest.Delete(ctx, u.ID)
	case entity.UserTypeFacilityUser:
		err = s.db.FacilityUser.Delete(ctx, u.ID)
	}
	return internalError(err)
}
