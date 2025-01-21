package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
)

// UserStatus - 購入者の状態
type UserStatus int32

const (
	UserStatusUnknown     UserStatus = 0
	UserStatusGuest       UserStatus = 1 // 未登録
	UserStatusProvisional UserStatus = 2 // 仮登録
	UserStatusVerified    UserStatus = 3 // 認証済み
	UserStatusDeactivated UserStatus = 4 // 無効
)

type User struct {
	response.User
	address Address
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
	case uentity.UserStatusDeactivated:
		return UserStatusDeactivated
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
	if user.Registered {
		return newMemberUser(user, address)
	}
	return newGuestUser(user, address)
}

func newMemberUser(user *uentity.User, address *uentity.Address) *User {
	return &User{
		User: response.User{
			ID:            user.ID,
			Status:        NewUserStatus(user.Status).Response(),
			Registered:    user.Registered,
			Username:      user.Member.Username,
			AccountID:     user.Member.AccountID,
			Lastname:      user.Member.Lastname,
			Firstname:     user.Member.Firstname,
			LastnameKana:  user.Member.LastnameKana,
			FirstnameKana: user.Member.FirstnameKana,
			Email:         user.Member.Email,
			PhoneNumber:   user.Member.PhoneNumber,
			ThumbnailURL:  user.Member.ThumbnailURL,
			CreatedAt:     jst.Unix(user.CreatedAt),
			UpdatedAt:     jst.Unix(user.UpdatedAt),
		},
		address: *NewAddress(address),
	}
}

func newGuestUser(user *uentity.User, address *uentity.Address) *User {
	return &User{
		User: response.User{
			ID:            user.ID,
			Status:        NewUserStatus(user.Status).Response(),
			Registered:    user.Registered,
			Username:      "ゲスト",
			AccountID:     "",
			Lastname:      user.Guest.Lastname,
			Firstname:     user.Guest.Firstname,
			LastnameKana:  user.Guest.LastnameKana,
			FirstnameKana: user.Guest.FirstnameKana,
			Email:         user.Guest.Email,
			PhoneNumber:   address.PhoneNumber, // 保持していないためアドレス帳の情報を使用
			CreatedAt:     user.CreatedAt.Unix(),
			UpdatedAt:     user.UpdatedAt.Unix(),
		},
		address: *NewAddress(address),
	}
}

func (u *User) Address() *Address {
	if u == nil || u.address.AddressID == "" {
		return nil
	}
	return &u.address
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

func (us Users) IDs() []string {
	return set.UniqBy(us, func(u *User) string {
		return u.ID
	})
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

func NewUserToList(user *User, order *sentity.AggregatedUserOrder) *UserToList {
	if order == nil {
		order = &sentity.AggregatedUserOrder{}
	}
	return &UserToList{
		UserToList: response.UserToList{
			ID:                user.ID,
			Lastname:          user.Lastname,
			Firstname:         user.Firstname,
			Email:             user.Email,
			Status:            user.Status,
			Registered:        user.Registered,
			PrefectureCode:    user.address.PrefectureCode,
			City:              user.address.City,
			PaymentTotalCount: order.OrderCount,
		},
	}
}

func (u *UserToList) Response() *response.UserToList {
	return &u.UserToList
}

func NewUsersToList(users Users, orders map[string]*sentity.AggregatedUserOrder) UsersToList {
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
