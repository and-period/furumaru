package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (s *service) MultiGetAdmins(ctx context.Context, in *user.MultiGetAdminsInput) (entity.Admins, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	admins, err := s.db.Admin.MultiGet(ctx, in.AdminIDs)
	return admins, exception.InternalError(err)
}

func (s *service) MultiGetAdminDevices(ctx context.Context, in *user.MultiGetAdminDevicesInput) ([]string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	auths, err := s.db.Admin.MultiGet(ctx, in.AdminIDs, "device")
	if err != nil {
		return nil, exception.InternalError(err)
	}
	return auths.Devices(), nil
}

func (s *service) GetAdmin(ctx context.Context, in *user.GetAdminInput) (*entity.Admin, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	admin, err := s.db.Admin.Get(ctx, in.AdminID)
	return admin, exception.InternalError(err)
}
