package service

import (
	"context"
	"fmt"

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

func (s *userService) ListAdmins(ctx context.Context, in *user.ListAdminsInput) (entity.Admins, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	roles, err := entity.NewAdminRoles(in.Roles)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrInvalidArgument, err.Error())
	}
	params := &database.ListAdminsParams{
		Roles:  roles,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	admins, err := s.db.Admin.List(ctx, params)
	return admins, exception.InternalError(err)
}

func (s *userService) MultiGetAdmins(ctx context.Context, in *user.MultiGetAdminsInput) (entity.Admins, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	admins, err := s.db.Admin.MultiGet(ctx, in.AdminIDs)
	return admins, exception.InternalError(err)
}

func (s *userService) GetAdmin(ctx context.Context, in *user.GetAdminInput) (*entity.Admin, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	a, err := s.db.Admin.Get(ctx, in.AdminID)
	return a, exception.InternalError(err)
}

func (s *userService) CreateAdmin(ctx context.Context, in *user.CreateAdminInput) (*entity.Admin, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	role, err := entity.NewAdminRole(in.Role)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrInvalidArgument, err.Error())
	}
	adminID := uuid.Base58Encode(uuid.New())
	admin := entity.NewAdmin(
		adminID, adminID,
		in.Lastname, in.Firstname,
		in.LastnameKana, in.FirstnameKana,
		in.Email, role,
	)
	if err := s.db.Admin.Create(ctx, admin); err != nil {
		return nil, exception.InternalError(err)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminCreateUserParams{
		Username: admin.CognitoID,
		Email:    admin.Email,
		Password: password,
	}
	if err := s.adminAuth.AdminCreateUser(ctx, params); err != nil {
		return nil, exception.InternalError(err)
	}
	s.logger.Debug("Create admin", zap.String("adminId", adminID), zap.String("password", password))
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		in := &messenger.NotifyRegisterAdminInput{
			Name:     admin.Name(),
			Email:    admin.Email,
			Password: password,
		}
		if err := s.messenger.NotifyRegisterAdmin(context.Background(), in); err != nil {
			s.logger.Debug("Notify register admin", zap.String("adminId", adminID), zap.Error(err))
		}
	}()
	return admin, nil
}

func (s *userService) UpdateAdminEmail(ctx context.Context, in *user.UpdateAdminEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return exception.InternalError(err)
	}
	a, err := s.db.Admin.GetByCognitoID(ctx, username, "id", "email")
	if err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		OldEmail:    a.Email,
		NewEmail:    in.Email,
	}
	err = s.adminAuth.ChangeEmail(ctx, params)
	return exception.InternalError(err)
}

func (s *userService) VerifyAdminEmail(ctx context.Context, in *user.VerifyAdminEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return exception.InternalError(err)
	}
	a, err := s.db.Admin.GetByCognitoID(ctx, username, "id")
	if err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		VerifyCode:  in.VerifyCode,
	}
	email, err := s.adminAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		return exception.InternalError(err)
	}
	err = s.db.Admin.UpdateEmail(ctx, a.ID, email)
	return exception.InternalError(err)
}

func (s *userService) UpdateAdminPassword(ctx context.Context, in *user.UpdateAdminPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: in.AccessToken,
		OldPassword: in.OldPassword,
		NewPassword: in.NewPassword,
	}
	err := s.adminAuth.ChangePassword(ctx, params)
	return exception.InternalError(err)
}
