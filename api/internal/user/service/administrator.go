package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/random"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListAdministrators(
	ctx context.Context, in *user.ListAdministratorsInput,
) (entity.Administrators, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	params := &database.ListAdministratorsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		administrators entity.Administrators
		total          int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		administrators, err = s.db.Administrator.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Administrator.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	return administrators, total, nil
}

func (s *service) MultiGetAdministrators(
	ctx context.Context, in *user.MultiGetAdministratorsInput,
) (entity.Administrators, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	administrators, err := s.db.Administrator.MultiGet(ctx, in.AdministratorIDs)
	return administrators, exception.InternalError(err)
}

func (s *service) GetAdministrator(
	ctx context.Context, in *user.GetAdministratorInput,
) (*entity.Administrator, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	administrator, err := s.db.Administrator.Get(ctx, in.AdministratorID)
	return administrator, exception.InternalError(err)
}

func (s *service) CreateAdministrator(
	ctx context.Context, in *user.CreateAdministratorInput,
) (*entity.Administrator, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	cognitoID := uuid.Base58Encode(uuid.New())
	password := random.NewStrings(size)
	if err := s.createCognitoAdmin(ctx, cognitoID, in.Email, password); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &entity.NewAdministratorParams{
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		Email:         in.Email,
		PhoneNumber:   in.PhoneNumber,
	}
	administrator := entity.NewAdministrator(params)
	auth := entity.NewAdminAuth(administrator.ID, cognitoID, entity.AdminRoleAdministrator)
	if err := s.db.Administrator.Create(ctx, auth, administrator); err != nil {
		return nil, exception.InternalError(err)
	}
	s.logger.Debug("Create administrator",
		zap.String("administratorId", administrator.ID), zap.String("password", password))
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyRegisterAdmin(context.Background(), administrator.ID, password)
		if err != nil {
			s.logger.Warn("Failed to notify register admin", zap.String("administratorId", administrator.ID), zap.Error(err))
		}
	}()
	return administrator, nil
}

func (s *service) createCognitoAdmin(ctx context.Context, cognitoID, email, password string) error {
	params := &cognito.AdminCreateUserParams{
		Username: cognitoID,
		Email:    email,
		Password: password,
	}
	return s.adminAuth.AdminCreateUser(ctx, params)
}

func (s *service) notifyRegisterAdmin(ctx context.Context, adminID, password string) error {
	in := &messenger.NotifyRegisterAdminInput{
		AdminID:  adminID,
		Password: password,
	}
	return s.messenger.NotifyRegisterAdmin(ctx, in)
}

func (s *service) notifyResetAdminPassword(ctx context.Context, adminID, password string) error {
	in := &messenger.NotifyResetAdminPasswordInput{
		AdminID:  adminID,
		Password: password,
	}
	return s.messenger.NotifyResetAdminPassword(ctx, in)
}
