package service

import (
	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/marche/api/internal/store/entity"
	uentity "github.com/and-period/marche/api/internal/user/entity"
)

type Staff struct {
	*response.Staff
}

type Staffs []*Staff

func NewStaff(staff *sentity.Staff, admin *uentity.Admin) *Staff {
	return &Staff{
		Staff: &response.Staff{
			ID:    admin.ID,
			Name:  admin.Name(),
			Email: admin.Email,
			Role:  int32(staff.Role),
		},
	}
}

func (s *Staff) Response() *response.Staff {
	return s.Staff
}

func NewStaffs(staffs sentity.Staffs, admins map[string]*uentity.Admin) Staffs {
	res := make(Staffs, len(staffs))
	for i, staff := range staffs {
		admin, ok := admins[staff.UserID]
		if !ok {
			continue
		}
		res[i] = NewStaff(staff, admin)
	}
	return res
}

func (ss Staffs) Response() []*response.Staff {
	res := make([]*response.Staff, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
