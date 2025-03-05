package entity

import (
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
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

// User - 購入者情報
type User struct {
	Member     `gorm:"-"`     // 会員情報
	Guest      `gorm:"-"`     // ゲスト情報
	ID         string         `gorm:"primaryKey;<-:create"` // ユーザーID
	Status     UserStatus     `gorm:"-"`                    // 購入者の状態
	Registered bool           `gorm:""`                     // 会員登録フラグ
	Device     string         `gorm:""`                     // デバイストークン(Push通知用)
	CreatedAt  time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time      `gorm:""`                     // 更新日時
	DeletedAt  gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Users []*User

type NewUserParams struct {
	Registered    bool
	CognitoID     string
	Username      string
	AccountID     string
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
	ProviderType  UserAuthProviderType
	Email         string
	PhoneNumber   string
}

func NewUser(params *NewUserParams) *User {
	var (
		member Member
		guest  Guest
	)
	userID := uuid.Base58Encode(uuid.New())
	if params.Registered {
		member.UserID = userID
		member.CognitoID = strings.ToLower(params.CognitoID) // Cognitoでは大文字小文字の区別がされず管理されているため
		member.Username = params.Username
		member.AccountID = params.AccountID
		member.Lastname = params.Lastname
		member.Firstname = params.Firstname
		member.LastnameKana = params.LastnameKana
		member.FirstnameKana = params.FirstnameKana
		member.ProviderType = params.ProviderType
		member.Email = params.Email
		member.PhoneNumber = params.PhoneNumber
	} else {
		guest.UserID = userID
		guest.Lastname = params.Lastname
		guest.Firstname = params.Firstname
		guest.LastnameKana = params.LastnameKana
		guest.FirstnameKana = params.FirstnameKana
		guest.Email = params.Email
	}
	return &User{
		ID:         userID,
		Registered: params.Registered,
		Member:     member,
		Guest:      guest,
	}
}

func (u *User) Name() string {
	if !u.Registered {
		return "ゲスト"
	}
	return u.Member.Name()
}

func (u *User) Username() string {
	if !u.Registered {
		return "名無しさん"
	}
	return u.Member.Username
}

func (u *User) Email() string {
	if u.Registered {
		return u.Member.Email
	}
	return u.Guest.Email
}

func (u *User) Fill(member *Member, guest *Guest) {
	u.Member = *member
	u.Guest = *guest
	u.SetStatus()
}

func (u *User) SetStatus() {
	if u == nil {
		return
	}
	switch {
	case !u.DeletedAt.Time.IsZero():
		u.Status = UserStatusDeactivated
	case !u.Registered:
		u.Status = UserStatusGuest
	case !u.VerifiedAt.IsZero():
		u.Status = UserStatusVerified
	default:
		u.Status = UserStatusProvisional
	}
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

func (us Users) GroupByRegistered() map[bool]Users {
	res := map[bool]Users{
		true:  make(Users, 0, len(us)),
		false: make(Users, 0, len(us)),
	}
	for _, u := range us {
		res[u.Registered] = append(res[u.Registered], u)
	}
	return res
}

func (us Users) Fill(members map[string]*Member, guests map[string]*Guest) {
	for _, u := range us {
		member, ok := members[u.ID]
		if !ok {
			member = &Member{}
		}
		guest, ok := guests[u.ID]
		if !ok {
			guest = &Guest{}
		}
		u.Fill(member, guest)
	}
}
