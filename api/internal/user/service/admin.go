package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"golang.org/x/sync/errgroup"
)

type Admin struct {
	response.Admin
}

type Admins []*Admin

func NewAdmin(admin *entity.Admin) *Admin {
	return &Admin{
		Admin: response.Admin{
			ID:            admin.ID,
			Role:          admin.Role,
			Lastname:      admin.Lastname,
			Firstname:     admin.Firstname,
			LastnameKana:  admin.LastnameKana,
			FirstnameKana: admin.FirstnameKana,
			Email:         admin.Email,
			CreatedAt:     admin.CreatedAt.Unix(),
			UpdatedAt:     admin.UpdatedAt.Unix(),
		},
	}
}

func (a *Admin) Response() *response.Admin {
	return &a.Admin
}

func NewAdmins(admins entity.Admins) Admins {
	res := make(Admins, len(admins))
	for i := range admins {
		res[i] = NewAdmin(admins[i])
	}
	return res
}

func (as Admins) Map() map[string]*Admin {
	res := make(map[string]*Admin, len(as))
	for _, a := range as {
		res[a.ID] = a
	}
	return res
}

func (as Admins) Response() []*response.Admin {
	res := make([]*response.Admin, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}

func (s *service) MultiGetAdmins(ctx context.Context, in *user.MultiGetAdminsInput) (entity.Admins, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	auths, err := s.db.AdminAuth.MultiGet(ctx, in.AdminIDs)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	res := make(entity.Admins, 0, len(auths))
	var m sync.Mutex
	eg, ectx := errgroup.WithContext(ctx)
	for role, as := range auths.GroupByRole() {
		role, adminIDs := role, as.AdminIDs()
		eg.Go(func() error {
			admins, err := s.getAdmins(ectx, adminIDs, role)
			if err != nil {
				return err
			}
			m.Lock()
			defer m.Unlock()
			res = append(res, admins...)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, exception.InternalError(err)
	}
	return res, nil
}

func (s *service) MultiGetAdminDevices(ctx context.Context, in *user.MultiGetAdminDevicesInput) ([]string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	auths, err := s.db.AdminAuth.MultiGet(ctx, in.AdminIDs, "device")
	if err != nil {
		return nil, exception.InternalError(err)
	}
	return auths.Devices(), nil
}

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

func (s *service) getAdmins(ctx context.Context, adminIDs []string, role entity.AdminRole) (entity.Admins, error) {
	switch role {
	case entity.AdminRoleAdministrator:
		administrators, err := s.db.Administrator.MultiGet(ctx, adminIDs)
		if err != nil {
			return nil, err
		}
		return entity.NewAdminsFromAdministrators(administrators), nil
	case entity.AdminRoleCoordinator:
		coordinators, err := s.db.Coordinator.MultiGet(ctx, adminIDs)
		if err != nil {
			return nil, err
		}
		return entity.NewAdminsFromCoordinators(coordinators), nil
	case entity.AdminRoleProducer:
		producers, err := s.db.Producer.MultiGet(ctx, adminIDs)
		if err != nil {
			return nil, err
		}
		return entity.NewAdminsFromProducers(producers), nil
	default:
		return nil, fmt.Errorf("api: invalid role: %w", exception.ErrInvalidArgument)
	}
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

func (a *Admin) Name() string {
	return strings.TrimSpace(strings.Join([]string{a.Lastname, a.Firstname}, " "))
}
