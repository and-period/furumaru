package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

// UserStatus - 購入者の状態
type UserStatus int32

const (
	UserStatusUnknown     UserStatus = 0
	UserStatusGuest       UserStatus = 1 // 未登録
	UserStatusProvisional UserStatus = 2 // 仮登録
	UserStatusVerified    UserStatus = 3 // 認証済み(初期設定前)
	UserStatusActivated   UserStatus = 4 // 認証済み(初期設定後)
)

type User struct {
	response.User
}

type Users []*User

type UserToList struct {
	response.UserToList
}

type UsersToList []*UserToList

func NewUserStatus(status uentity.UserStatus) UserStatus {
	switch status {
	case uentity.UserStatusGuest:
		return UserStatusGuest
	case uentity.UserStatusProvisional:
		return UserStatusProvisional
	case uentity.UserStatusVerified:
		return UserStatusVerified
	case uentity.UserStatusActivated:
		return UserStatusActivated
	default:
		return UserStatusUnknown
	}
}

func (s UserStatus) Response() int32 {
	return int32(s)
}

func NewUser(user *uentity.User, address *uentity.Address) *User {
	if address == nil {
		address = &uentity.Address{}
	}
	return &User{
		User: response.User{
			ID:         user.ID,
			Registered: user.Registered,
			Status:     NewUserStatus(user.Status).Response(),
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

func NewUserToList(user *User, order *sentity.AggregatedOrder) *UserToList {
	if order == nil {
		order = &sentity.AggregatedOrder{}
	}
	return &UserToList{
		UserToList: response.UserToList{
			ID:             user.ID,
			Lastname:       user.Lastname,
			Firstname:      user.Firstname,
			Email:          user.Email,
			Registered:     user.Registered,
			PrefectureCode: user.PrefectureCode,
			City:           user.City,
			TotalOrder:     order.OrderCount,
			TotalAmount:    order.Subtotal,
		},
	}
}

func (u *UserToList) Response() *response.UserToList {
	return &u.UserToList
}

func NewUsersToList(users Users, orders map[string]*sentity.AggregatedOrder) UsersToList {
	res := make(UsersToList, len(users))
	for i := range users {
		res[i] = NewUserToList(users[i], orders[users[i].ID])
	}
	return res
}

func (us UsersToList) Response() []*response.UserToList {
	res := make([]*response.UserToList, len(us))
	for i := range us {
		res[i] = us[i].Response()
	}
	return res
}
