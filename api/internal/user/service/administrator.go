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
)

func (s *userService) ListAdministrators(
	ctx context.Context, in *user.ListAdministratorsInput,
) (entity.Administrators, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &database.ListAdministratorsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	administrators, err := s.db.Administrator.List(ctx, params)
	return administrators, exception.InternalError(err)
}

func (s *userService) GetAdministrator(
	ctx context.Context, in *user.GetAdministratorInput,
) (*entity.Administrator, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	administrator, err := s.db.Administrator.Get(ctx, in.AdministratorID)
	return administrator, exception.InternalError(err)
}

func (s *userService) CreateAdministrator(
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
	administratorID := uuid.Base58Encode(uuid.New())
	params := &entity.NewAdministratorParams{
		ID:            administratorID,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.Lastname,
		FirstnameKana: in.Firstname,
		Email:         in.Email,
		PhoneNumber:   in.PhoneNumber,
	}
	auth := entity.NewAdminAuth(administratorID, cognitoID, entity.AdminRoleAdministrator)
	administrator := entity.NewAdministrator(params)
	if err := s.db.Administrator.Create(ctx, auth, administrator); err != nil {
		return nil, exception.InternalError(err)
	}
	s.logger.Debug("Create administrator",
		zap.String("administratorId", administratorID), zap.String("password", password))
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		err := s.notifyRegisterAdmin(context.Background(), administrator.Name(), administrator.Email, password)
		if err != nil {
			s.logger.Warn("Failed to notify register admin", zap.String("administratorId", administratorID), zap.Error(err))
		}
	}()
	return administrator, nil
}

func (s *userService) createCognitoAdmin(ctx context.Context, cognitoID, email, password string) error {
	params := &cognito.AdminCreateUserParams{
		Username: cognitoID,
		Email:    email,
		Password: password,
	}
	return s.adminAuth.AdminCreateUser(ctx, params)
}

func (s *userService) notifyRegisterAdmin(ctx context.Context, name, email, password string) error {
	in := &messenger.NotifyRegisterAdminInput{
		Name:     name,
		Email:    email,
		Password: password,
	}
	return s.messenger.NotifyRegisterAdmin(ctx, in)
}
