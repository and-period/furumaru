package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
)

func (s *service) MultiGetAdmins(
	ctx context.Context,
	in *user.MultiGetAdminsInput,
) (entity.Admins, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	admins, err := s.db.Admin.MultiGet(ctx, in.AdminIDs)
	return admins, internalError(err)
}

func (s *service) MultiGetAdminDevices(
	ctx context.Context,
	in *user.MultiGetAdminDevicesInput,
) ([]string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	auths, err := s.db.Admin.MultiGet(ctx, in.AdminIDs, "device")
	if err != nil {
		return nil, internalError(err)
	}
	return auths.Devices(), nil
}

func (s *service) GetAdmin(ctx context.Context, in *user.GetAdminInput) (*entity.Admin, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	admin, err := s.db.Admin.Get(ctx, in.AdminID)
	return admin, internalError(err)
}

func (s *service) ForgotAdminPassword(
	ctx context.Context,
	in *user.ForgotAdminPasswordInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	admin, err := s.db.Admin.GetByEmail(ctx, in.Email, "cognito_id")
	if err != nil {
		return internalError(err)
	}
	if err := s.adminAuth.ForgotPassword(ctx, admin.CognitoID); err != nil {
		return fmt.Errorf("%w: %s", exception.ErrNotFound, err.Error())
	}
	return nil
}

func (s *service) VerifyAdminPassword(
	ctx context.Context,
	in *user.VerifyAdminPasswordInput,
) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	admin, err := s.db.Admin.GetByEmail(ctx, in.Email, "cognito_id")
	if err != nil {
		return internalError(err)
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    admin.CognitoID,
		VerifyCode:  in.VerifyCode,
		NewPassword: in.NewPassword,
	}
	err = s.adminAuth.ConfirmForgotPassword(ctx, params)
	return internalError(err)
}
