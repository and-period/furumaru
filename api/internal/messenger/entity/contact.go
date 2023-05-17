package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// お問い合わせステータス
type ContactStatus int32

const (
	ContactStatusUnknown    ContactStatus = 0
	ContactStatusWaiting    ContactStatus = 1 // 未着手
	ContactStatusInprogress ContactStatus = 2 // 進行中
	ContactStatusDone       ContactStatus = 3 // 完了
	ContactStatusDiscard    ContactStatus = 4 // 対応不要
)

// Contact - お問い合わせ情報
type Contact struct {
	ID          string         `gorm:"primaryKey;<-:create"` // お問い合わせID
	Title       string         `gorm:""`                     // 件名
	CategoryID  string         `gorm:"default:null"`         // 問い合わせ種別ID
	Content     string         `gorm:""`                     // 内容
	Username    string         `gorm:""`                     // 氏名
	UserID      string         `gorm:"default:null"`         // ユーザーID
	Email       string         `gorm:""`                     // メールアドレス
	PhoneNumber string         `gorm:""`                     // 電話番号
	Status      ContactStatus  `gorm:""`                     // ステータス
	ResponderID string         `gorm:"default:null"`         // 対応者ID
	Note        string         `gorm:""`                     // 対応者メモ
	CreatedAt   time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time      `gorm:""`                     // 更新日時
	DeletedAt   gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Contacts []*Contact

type NewContactParams struct {
	Title       string
	Content     string
	Username    string
	Email       string
	PhoneNumber string
	Note        string
}

func NewContact(params *NewContactParams) *Contact {
	return &Contact{
		ID:          uuid.Base58Encode(uuid.New()),
		Title:       params.Title,
		Content:     params.Content,
		Username:    params.Username,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		Status:      ContactStatusUnknown,
		Note:        "",
	}
}

func (c *Contact) Fill(categoryId, userId, responderId string) {
	c.CategoryID = categoryId
	c.UserID = userId
	c.ResponderID = responderId
}
