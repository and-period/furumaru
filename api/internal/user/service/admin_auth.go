package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
)

func (s *service) SignInAdmin(ctx context.Context, in *user.SignInAdminInput) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	rs, err := s.adminAuth.SignIn(ctx, in.Key, in.Password)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	auth, err := s.getAdminAuth(ctx, rs)
	return auth, exception.InternalError(err)
}

func (s *service) SignOutAdmin(ctx context.Context, in *user.SignOutAdminInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.adminAuth.SignOut(ctx, in.AccessToken)
	return exception.InternalError(err)
}

func (s *service) GetAdminAuth(ctx context.Context, in *user.GetAdminAuthInput) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	rs := &cognito.AuthResult{AccessToken: in.AccessToken}
	auth, err := s.getAdminAuth(ctx, rs)
	return auth, exception.InternalError(err)
}

func (s *service) RefreshAdminToken(
	ctx context.Context, in *user.RefreshAdminTokenInput,
) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	rs, err := s.adminAuth.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	auth, err := s.getAdminAuth(ctx, rs)
	return auth, exception.InternalError(err)
}

func (s *service) UpdateAdminEmail(ctx context.Context, in *user.UpdateAdminEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return exception.InternalError(err)
	}
	auth, err := s.db.AdminAuth.GetByCognitoID(ctx, username, "admin_id", "role")
	if err != nil {
		return exception.InternalError(err)
	}
	admin, err := s.getAdmin(ctx, auth.AdminID, auth.Role)
	if err != nil {
		return exception.InternalError(err)
	}
	if admin.Email == in.Email {
		return fmt.Errorf("this admin does not need to be changed email: %w", exception.ErrFailedPrecondition)
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		OldEmail:    admin.Email,
		NewEmail:    in.Email,
	}
	err = s.adminAuth.ChangeEmail(ctx, params)
	return exception.InternalError(err)
}

func (s *service) VerifyAdminEmail(ctx context.Context, in *user.VerifyAdminEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return exception.InternalError(err)
	}
	auth, err := s.db.AdminAuth.GetByCognitoID(ctx, username, "admin_id", "role")
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
	switch auth.Role {
	case entity.AdminRoleAdministrator:
		err = s.db.Administrator.UpdateEmail(ctx, auth.AdminID, email)
	case entity.AdminRoleCoordinator:
		err = s.db.Coordinator.UpdateEmail(ctx, auth.AdminID, email)
	case entity.AdminRoleProducer:
		err = s.db.Producer.UpdateEmail(ctx, auth.AdminID, email)
	}
	return exception.InternalError(err)
}

func (s *service) UpdateAdminPassword(ctx context.Context, in *user.UpdateAdminPasswordInput) error {
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

func (s *service) getAdminAuth(ctx context.Context, rs *cognito.AuthResult) (*entity.AdminAuth, error) {
	username, err := s.adminAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, err
	}
	auth, err := s.db.AdminAuth.GetByCognitoID(ctx, username)
	if err != nil {
		return nil, err
	}
	auth.Fill(rs)
	return auth, nil
}
