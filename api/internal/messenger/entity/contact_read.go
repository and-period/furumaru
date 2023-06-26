package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

type ContactUserType int32

const (
	ContactUserTypeUnknown ContactUserType = iota // 不明
	ContactUserTypeAdmin                          // 管理者
	ContactUserTypeUser                           // ユーザー
	ContactUserTypeGuest                          // ゲスト(ユーザIDなし)
)

type ContactRead struct {
	ID        string          `gorm:"primaryKey;<-:create"` // お問い合わせ会話履歴ID
	ContactID string          `gorm:""`                     // お問い合わせID
	UserID    string          `gorm:"default:null"`         // 送信者ID(ゲストの場合null)
	UserType  ContactUserType `gorm:""`                     // 送信者の種別
	ReadFlag  bool            `gorm:""`                     // 既読フラグ
	CreatedAt time.Time       `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time       `gorm:""`                     // 更新日時
}

type ContactReads []*ContactRead

type NewContactReadParams struct {
	ContactID string
	UserType  ContactUserType
	ReadFlag  bool
}

func NewContactRead(params *NewContactReadParams) *ContactRead {
	return &ContactRead{
		ID:        uuid.Base58Encode(uuid.New()),
		ContactID: params.ContactID,
		UserType:  params.UserType,
		ReadFlag:  params.ReadFlag,
	}
}

func (c *ContactRead) Fill(userID string) {
	if userID != "" {
		c.UserID = userID
	}
}
