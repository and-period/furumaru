package service

import (
	"strings"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type User struct {
	response.User
}

type Users []*User

func NewUser(user *entity.User) *User {
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
			Prefecture:    user.Customer.Prefecture,
			City:          user.Customer.City,
			AddressLine1:  user.Customer.AddressLine1,
			AddressLine2:  user.Customer.AddressLine2,
			CreatedAt:     user.Customer.CreatedAt.Unix(),
			UpdatedAt:     user.Customer.UpdatedAt.Unix(),
		},
	}
}

func (u *User) Name() string {
	return strings.TrimSpace(strings.Join([]string{u.Lastname, u.Firstname}, " "))
}

func (u *User) Response() *response.User {
	return &u.User
}

func NewUsers(users entity.Users) Users {
	res := make(Users, len(users))
	for i := range users {
		res[i] = NewUser(users[i])
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
