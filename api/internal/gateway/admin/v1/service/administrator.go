package service

import (
	"fmt"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Administrator struct {
	types.Administrator
}

type Administrators []*Administrator

func NewAdministrator(admin *entity.Administrator) *Administrator {
	return &Administrator{
		Administrator: types.Administrator{
			ID:            admin.ID,
			Status:        NewAdminStatus(admin.Status).Response(),
			Lastname:      admin.Lastname,
			Firstname:     admin.Firstname,
			LastnameKana:  admin.LastnameKana,
			FirstnameKana: admin.FirstnameKana,
			Email:         admin.Email,
			PhoneNumber:   admin.PhoneNumber,
			CreatedAt:     admin.CreatedAt.Unix(),
			UpdatedAt:     admin.CreatedAt.Unix(),
		},
	}
}

func (a *Administrator) AuthUser() *AuthUser {
	return &AuthUser{
		AuthUser: types.AuthUser{
			AdminID:  a.ID,
			ShopIDs:  []string{},
			Type:     types.AdminTypeAdministrator,
			Username: fmt.Sprintf("%s %s", a.Lastname, a.Firstname),
			Email:    a.Email,
		},
	}
}

func (a *Administrator) Response() *types.Administrator {
	return &a.Administrator
}

func NewAdministrators(admins entity.Administrators) Administrators {
	res := make(Administrators, len(admins))
	for i := range admins {
		res[i] = NewAdministrator(admins[i])
	}
	return res
}

func (as Administrators) Response() []*types.Administrator {
	res := make([]*types.Administrator, len(as))
	for i := range as {
		res[i] = as[i].Response()
	}
	return res
}

func (as Administrators) Map() map[string]*Administrator {
	res := make(map[string]*Administrator, len(as))
	for _, a := range as {
		res[a.ID] = a
	}
	return res
}
