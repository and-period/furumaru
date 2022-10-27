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
	adminParams := &entity.NewAdminParams{
		CognitoID:     cognitoID,
		Role:          entity.AdminRoleAdministrator,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		Email:         in.Email,
	}
	params := &entity.NewAdministratorParams{
		Admin:       entity.NewAdmin(adminParams),
		PhoneNumber: in.PhoneNumber,
	}
	administrator := entity.NewAdministrator(params)
	auth := s.createCognitoAdmin(cognitoID, in.Email, password)
	if err := s.db.Administrator.Create(ctx, administrator, auth); err != nil {
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

func (s *service) UpdateAdministrator(ctx context.Context, in *user.UpdateAdministratorInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	params := &database.UpdateAdministratorParams{
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		PhoneNumber:   in.PhoneNumber,
	}
	err := s.db.Administrator.Update(ctx, in.AdministratorID, params)
	return exception.InternalError(err)
}

func (s *service) UpdateAdministratorEmail(ctx context.Context, in *user.UpdateAdministratorEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	administrator, err := s.db.Administrator.Get(ctx, in.AdministratorID)
	if err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.AdminChangeEmailParams{
		Username: administrator.CognitoID,
		Email:    in.Email,
	}
	if err := s.adminAuth.AdminChangeEmail(ctx, params); err != nil {
		return exception.InternalError(err)
	}
	err = s.db.Admin.UpdateEmail(ctx, in.AdministratorID, in.Email)
	return exception.InternalError(err)
}

func (s *service) ResetAdministratorPassword(ctx context.Context, in *user.ResetAdministratorPasswordInput) error {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	administrator, err := s.db.Administrator.Get(ctx, in.AdministratorID)
	if err != nil {
		return exception.InternalError(err)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminChangePasswordParams{
		Username:  administrator.CognitoID,
		Password:  password,
		Permanent: true,
	}
	if err := s.adminAuth.AdminChangePassword(ctx, params); err != nil {
		return exception.InternalError(err)
	}
	s.logger.Debug("Reset administrator password",
		zap.String("administrator", in.AdministratorID), zap.String("password", password),
	)
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyResetAdminPassword(context.Background(), in.AdministratorID, password)
		if err != nil {
			s.logger.Warn("Failed to notify reset admin password",
				zap.String("administrator", in.AdministratorID), zap.Error(err),
			)
		}
	}()
	return nil
}

func (s *service) DeleteAdministrator(ctx context.Context, in *user.DeleteAdministratorInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Administrator.Delete(ctx, in.AdministratorID, s.deleteCognitoAdmin(in.AdministratorID))
	return exception.InternalError(err)
}

func (s *service) createCognitoAdmin(cognitoID, email, password string) func(context.Context) error {
	params := &cognito.AdminCreateUserParams{
		Username: cognitoID,
		Email:    email,
		Password: password,
	}
	return func(ctx context.Context) error {
		return s.adminAuth.AdminCreateUser(ctx, params)
	}
}

func (s *service) deleteCognitoAdmin(cognitoID string) func(context.Context) error {
	return func(ctx context.Context) error {
		return s.adminAuth.DeleteUser(ctx, cognitoID)
	}
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
