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

func NewStaff(staff *sentity.Staff, user *uentity.Shop) *Staff {
	return &Staff{
		Staff: &response.Staff{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  int32(staff.Role),
		},
	}
}

func (s *Staff) Response() *response.Staff {
	return s.Staff
}

func NewStaffs(staffs sentity.Staffs, users map[string]*uentity.Shop) Staffs {
	res := make(Staffs, len(staffs))
	for i, staff := range staffs {
		user, ok := users[staff.UserID]
		if !ok {
			continue
		}
		res[i] = NewStaff(staff, user)
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
