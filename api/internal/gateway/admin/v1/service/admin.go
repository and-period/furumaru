package service

import (
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
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

func (a *Admin) Name() string {
	return strings.TrimSpace(strings.Join([]string{a.Lastname, a.Firstname}, " "))
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
