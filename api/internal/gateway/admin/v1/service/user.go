package service

import (
	"strings"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

type User struct {
	response.User
}

type Users []*User

type UserList struct {
	response.UserList
}

type UserLists []*UserList

func NewUser(user *uentity.User) *User {
	return &User{
		User: response.User{
			ID:            user.ID,
			Lastname:      user.Customer.Lastname,
			Firstname:     user.Customer.Firstname,
			LastnameKana:  user.Customer.LastnameKana,
			FirstnameKana: user.Customer.FirstnameKana,
			Registered:    user.Registered,
			Email:         user.Email(),
			PhoneNumber:   user.PhoneNumber(),
			PostalCode:    user.Customer.PostalCode,
			Prefecture:    codes.PrefectureNames[user.Customer.Prefecture],
			City:          user.Customer.City,
			AddressLine1:  user.Customer.AddressLine1,
			AddressLine2:  user.Customer.AddressLine2,
			CreatedAt:     user.CreatedAt.Unix(),
			UpdatedAt:     user.UpdatedAt.Unix(),
		},
	}
}

func (u *User) Name() string {
	return strings.TrimSpace(strings.Join([]string{u.Lastname, u.Firstname}, " "))
}

func (u *User) Response() *response.User {
	return &u.User
}

func NewUsers(users uentity.Users) Users {
	res := make(Users, len(users))
	for i := range users {
		res[i] = NewUser(users[i])
	}
	return res
}

func (us Users) IDs() []string {
	res := make([]string, len(us))
	for i := range us {
		res[i] = us[i].ID
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

func NewUserList(user *User, order *sentity.AggregatedOrder) *UserList {
	if order == nil {
		order = &sentity.AggregatedOrder{}
	}
	return &UserList{
		UserList: response.UserList{
			ID:          user.ID,
			Lastname:    user.Lastname,
			Firstname:   user.Firstname,
			Registered:  user.Registered,
			Prefecture:  user.Prefecture,
			City:        user.City,
			TotalOrder:  order.OrderCount,
			TotalAmount: order.Subtotal,
		},
	}
}

func (u *UserList) Response() *response.UserList {
	return &u.UserList
}

func NewUserLists(users Users, orders map[string]*sentity.AggregatedOrder) UserLists {
	res := make(UserLists, len(users))
	for i := range users {
		res[i] = NewUserList(users[i], orders[users[i].ID])
	}
	return res
}

func (us UserLists) Response() []*response.UserList {
	res := make([]*response.UserList, len(us))
	for i := range us {
		res[i] = us[i].Response()
	}
	return res
}
