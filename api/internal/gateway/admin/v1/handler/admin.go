package handler

import (
	"context"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/user"
)

func (h *handler) multiGetAdmins(ctx context.Context, adminIDs []string) (service.Admins, error) {
	if len(adminIDs) == 0 {
		return service.Admins{}, nil
	}
	in := &user.MultiGetAdminsInput{
		AdminIDs: adminIDs,
	}
	admins, err := h.user.MultiGetAdmins(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewAdmins(admins), nil
}

func (h *handler) getAdmin(ctx context.Context, adminID string) (*service.Admin, error) {
	in := &user.GetAdminInput{
		AdminID: adminID,
	}
	admin, err := h.user.GetAdmin(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewAdmin(admin), nil
}
