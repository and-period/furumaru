package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// お問い合わせステータス
type ContactStatus int32

const (
	ContactStatusUnknown    ContactStatus = 0
	ContactStatusToDo       ContactStatus = 1 // ToDo
	ContactStatusInprogress ContactStatus = 2 // 進行中
	ContactStatusDone       ContactStatus = 3 // 完了
	ContactStatusDiscard    ContactStatus = 4 // 対応不要
)

// お問い合わせ優先度
type ContactPriority int32

const (
	ContactPriorityUnknown ContactPriority = 0
	ContactPriorityLow     ContactPriority = 1 // 優先度・低
	ContactPriorityMiddle  ContactPriority = 2 // 優先度・中
	ContactPriorityHigh    ContactPriority = 3 // 優先度・高
)

// Contact - お問い合わせ情報
type Contact struct {
	ID          string          `gorm:"primaryKey;<-:create"` // お問い合わせID
	Title       string          `gorm:""`                     // 件名
	Content     string          `gorm:""`                     // 内容
	Username    string          `gorm:""`                     // 氏名
	Email       string          `gorm:""`                     // メールアドレス
	PhoneNumber string          `gorm:""`                     // 電話番号
	Status      ContactStatus   `gorm:""`                     // ステータス
	Priority    ContactPriority `gorm:""`                     // 対応優先度
	Note        string          `gorm:""`                     // 対応者メモ
	CreatedAt   time.Time       `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time       `gorm:""`                     // 更新日時
}

func NewContact(title, content, username, email, phoneNumber string) *Contact {
	return &Contact{
		ID:          uuid.Base58Encode(uuid.New()),
		Title:       title,
		Content:     content,
		Username:    username,
		Email:       email,
		PhoneNumber: phoneNumber,
		Status:      ContactStatusUnknown,
		Priority:    ContactPriorityUnknown,
		Note:        "",
	}
}
