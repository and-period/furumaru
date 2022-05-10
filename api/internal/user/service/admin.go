package service

import (
	"context"
	"fmt"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/random"
	"github.com/and-period/marche/api/pkg/uuid"
)

func (s *userService) ListAdmins(ctx context.Context, in *ListAdminsInput) (entity.Admins, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	return nil, ErrNotImplemented
}

func (s *userService) MultiGetAdmins(ctx context.Context, in *MultiGetAdminsInput) (entity.Admins, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	return nil, ErrNotImplemented
}

func (s *userService) GetAdmin(ctx context.Context, in *GetAdminInput) (*entity.Admin, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	a, err := s.db.Admin.Get(ctx, in.AdminID)
	return a, userError(err)
}

func (s *userService) CreateAdmin(ctx context.Context, in *CreateAdminInput) (*entity.Admin, error) {
	const size = 8
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	if err := in.Role.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	}
	adminID := uuid.Base58Encode(uuid.New())
	admin := entity.NewAdmin(
		adminID, adminID,
		in.Lastname, in.Firstname,
		in.LastnameKana, in.FirstnameKana,
		in.Email, in.Role,
	)
	if err := s.db.Admin.Create(ctx, admin); err != nil {
		return nil, userError(err)
	}
	password := random.NewStrings(size)
	params := &cognito.AdminCreateUserParams{
		Username: admin.CognitoID,
		Email:    admin.Email,
		Password: password,
	}
	if err := s.adminAuth.AdminCreateUser(ctx, params); err != nil {
		return nil, userError(err)
	}
	// TODO: 管理者登録通知を送信
	return admin, nil
}

func (s *userService) UpdateAdminEmail(ctx context.Context, in *UpdateAdminEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return userError(err)
	}
	a, err := s.db.Admin.GetByCognitoID(ctx, username, "id", "email")
	if err != nil {
		return userError(err)
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		OldEmail:    a.Email,
		NewEmail:    in.Email,
	}
	err = s.adminAuth.ChangeEmail(ctx, params)
	return userError(err)
}

func (s *userService) VerifyAdminEmail(ctx context.Context, in *VerifyAdminEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return userError(err)
	}
	a, err := s.db.Admin.GetByCognitoID(ctx, username, "id")
	if err != nil {
		return userError(err)
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		VerifyCode:  in.VerifyCode,
	}
	email, err := s.adminAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		return userError(err)
	}
	err = s.db.Admin.UpdateEmail(ctx, a.ID, email)
	return userError(err)
}

func (s *userService) UpdateAdminPassword(ctx context.Context, in *UpdateAdminPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: in.AccessToken,
		OldPassword: in.OldPassword,
		NewPassword: in.NewPassword,
	}
	err := s.adminAuth.ChangePassword(ctx, params)
	return userError(err)
}
