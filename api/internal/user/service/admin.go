package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func (s *service) GetAdmin(ctx context.Context, in *user.GetAdminInput) (*entity.Admin, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	auth, err := s.db.AdminAuth.GetByAdminID(ctx, in.AdminID, "role")
	if err != nil {
		return nil, exception.InternalError(err)
	}
	admin, err := s.getAdmin(ctx, in.AdminID, auth.Role)
	return admin, exception.InternalError(err)
}

func (s *service) getAdmin(ctx context.Context, adminID string, role entity.AdminRole) (*entity.Admin, error) {
	switch role {
	case entity.AdminRoleAdministrator:
		administrator, err := s.db.Administrator.Get(ctx, adminID)
		if err != nil {
			return nil, err
		}
		return entity.NewAdminFromAdministrator(administrator), nil
	case entity.AdminRoleCoordinator:
		coordinator, err := s.db.Coordinator.Get(ctx, adminID)
		if err != nil {
			return nil, err
		}
		return entity.NewAdminFromCoordinator(coordinator), nil
	case entity.AdminRoleProducer:
		producer, err := s.db.Producer.Get(ctx, adminID)
		if err != nil {
			return nil, err
		}
		return entity.NewAdminFromProducer(producer), nil
	default:
		return nil, fmt.Errorf("api: invalid role: %w", exception.ErrInvalidArgument)
	}
}
