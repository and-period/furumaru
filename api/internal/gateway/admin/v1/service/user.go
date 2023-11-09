package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

type User struct {
	response.User
}

type Users []*User

type UserSummary struct {
	response.UserSummary
}

type UserSummaries []*UserSummary

func NewUser(user *uentity.User, address *uentity.Address) *User {
	if address == nil {
		address = &uentity.Address{}
	}
	return &User{
		User: response.User{
			ID:         user.ID,
			Registered: user.Registered,
			Email:      user.Email(),
			Address:    NewAddress(address).Response(),
			CreatedAt:  user.CreatedAt.Unix(),
			UpdatedAt:  user.UpdatedAt.Unix(),
		},
	}
}

func (u *User) Response() *response.User {
	return &u.User
}

func NewUsers(users uentity.Users, addresses map[string]*uentity.Address) Users {
	res := make(Users, len(users))
	for i, u := range users {
		res[i] = NewUser(u, addresses[u.ID])
	}
	return res
}

func (us Users) Map() map[string]*User {
	res := make(map[string]*User, len(us))
	for _, u := range us {
		res[u.ID] = u
	}
	return res
}

func (us Users) Response() []*response.User {
	res := make([]*response.User, len(us))
	for i := range us {
		res[i] = us[i].Response()
	}
	return res
}

func NewUserSummary(user *User, order *sentity.AggregatedOrder) *UserSummary {
	if order == nil {
		order = &sentity.AggregatedOrder{}
	}
	return &UserSummary{
		UserSummary: response.UserSummary{
			ID:             user.ID,
			Lastname:       user.Lastname,
			Firstname:      user.Firstname,
			Registered:     user.Registered,
			PrefectureCode: user.PrefectureCode,
			City:           user.City,
			TotalOrder:     order.OrderCount,
			TotalAmount:    order.Subtotal,
		},
	}
}

func (u *UserSummary) Response() *response.UserSummary {
	return &u.UserSummary
}

func NewUserSummaries(users Users, orders map[string]*sentity.AggregatedOrder) UserSummaries {
	res := make(UserSummaries, len(users))
	for i := range users {
		res[i] = NewUserSummary(users[i], orders[users[i].ID])
	}
	return res
}

func (us UserSummaries) Response() []*response.UserSummary {
	res := make([]*response.UserSummary, len(us))
	for i := range us {
		res[i] = us[i].Response()
	}
	return res
}
