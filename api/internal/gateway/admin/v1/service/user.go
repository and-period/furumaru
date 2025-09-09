package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
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
	types.User
	address Address
}

type Users []*User

type UserToList struct {
	types.UserToList
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
	switch user.Type {
	case uentity.UserTypeMember:
		return newMemberUser(user, address)
	case uentity.UserTypeGuest:
		return newGuestUser(user, address)
	case uentity.UserTypeFacilityUser:
		return newFacilityUser(user, address)
	default:
		return nil
	}
}

func newMemberUser(user *uentity.User, address *uentity.Address) *User {
	return &User{
		User: types.User{
			ID:            user.ID,
			Status:        NewUserStatus(user.Status).Response(),
			Registered:    user.Registered,
			Username:      user.Member.Username,
			AccountID:     user.AccountID,
			Lastname:      user.Member.Lastname,
			Firstname:     user.Member.Firstname,
			LastnameKana:  user.Member.LastnameKana,
			FirstnameKana: user.Member.FirstnameKana,
			Email:         user.Member.Email,
			PhoneNumber:   user.Member.PhoneNumber,
			ThumbnailURL:  user.ThumbnailURL,
			CreatedAt:     jst.Unix(user.CreatedAt),
			UpdatedAt:     jst.Unix(user.UpdatedAt),
		},
		address: *NewAddress(address),
	}
}

func newGuestUser(user *uentity.User, address *uentity.Address) *User {
	return &User{
		User: types.User{
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
			CreatedAt:     jst.Unix(user.CreatedAt),
			UpdatedAt:     jst.Unix(user.UpdatedAt),
		},
		address: *NewAddress(address),
	}
}

func newFacilityUser(user *uentity.User, _ *uentity.Address) *User {
	return &User{
		User: types.User{
			ID:            user.ID,
			Status:        NewUserStatus(user.Status).Response(),
			Registered:    user.Registered,
			Username:      "外部宿泊施設利用者",
			AccountID:     "",
			Lastname:      user.FacilityUser.Lastname,
			Firstname:     user.FacilityUser.Firstname,
			LastnameKana:  user.FacilityUser.LastnameKana,
			FirstnameKana: user.FacilityUser.FirstnameKana,
			Email:         user.FacilityUser.Email,
			PhoneNumber:   user.FacilityUser.PhoneNumber,
			CreatedAt:     jst.Unix(user.CreatedAt),
			UpdatedAt:     jst.Unix(user.UpdatedAt),
		},
		address: Address{},
	}
}

func (u *User) Address() *Address {
	if u == nil || u.address.AddressID == "" {
		return nil
	}
	return &u.address
}

func (u *User) Response() *types.User {
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

func (us Users) Response() []*types.User {
	res := make([]*types.User, len(us))
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
		UserToList: types.UserToList{
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

func (u *UserToList) Response() *types.UserToList {
	return &u.UserToList
}

func NewUsersToList(users Users, orders map[string]*sentity.AggregatedUserOrder) UsersToList {
	res := make(UsersToList, len(users))
	for i := range users {
		res[i] = NewUserToList(users[i], orders[users[i].ID])
	}
	return res
}

func (us UsersToList) Response() []*types.UserToList {
	res := make([]*types.UserToList, len(us))
	for i := range us {
		res[i] = us[i].Response()
	}
	return res
}
