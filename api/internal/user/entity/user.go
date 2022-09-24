package entity

import (
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// User - 購入者情報
type User struct {
	Member     `gorm:"-"` // 会員情報
	Guest      `gorm:"-"` // ゲスト情報
	Customer   `gorm:"-"` // 購入者情報
	ID         string     `gorm:"primaryKey;<-:create"` // ユーザーID
	Registered bool       `gorm:""`                     // 会員登録フラグ
	Device     string     `gorm:""`                     // デバイストークン(Push通知用)
	CreatedAt  time.Time  `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time  `gorm:""`                     // 更新日時
}

type Users []*User

type NewUserParams struct {
	Registered   bool
	CognitoID    string
	ProviderType ProviderType
	Email        string
	PhoneNumber  string
}

func NewUser(params *NewUserParams) *User {
	var (
		member Member
		guest  Guest
	)
	userID := uuid.Base58Encode(uuid.New())
	if params.Registered {
		member.UserID = userID
		member.CognitoID = params.CognitoID
		member.ProviderType = params.ProviderType
		member.Email = params.Email
		member.PhoneNumber = params.PhoneNumber
	} else {
		guest.UserID = userID
		guest.Email = params.Email
		guest.PhoneNumber = params.PhoneNumber
	}
	return &User{
		ID:         userID,
		Registered: params.Registered,
		Member:     member,
		Guest:      guest,
	}
}

func (u *User) Name() string {
	return strings.TrimSpace(strings.Join([]string{u.Lastname, u.Firstname}, " "))
}

func (u *User) Fill(customer *Customer, member *Member, guest *Guest) {
	u.Customer = *customer
	u.Member = *member
	u.Guest = *guest
}

func (us Users) IDs() []string {
	res := make([]string, len(us))
	for i := range us {
		res[i] = us[i].ID
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
