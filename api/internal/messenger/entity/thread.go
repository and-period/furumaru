package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// Thread - お問い合わせ会話履歴
type Thread struct {
	ID        string         `gorm:"primaryKey;<-:create"` // お問い合わせ会話履歴ID
	ContactID string         `gorm:""`                     // お問い合わせID
	UserID    string         `gorm:"default:null"`         // 送信者ID(ゲストの場合null)
	UserType  int32          `gorm:""`                     // 送信者の種別(不明:0, admin:1, uer:2, guest:3)
	Content   string         `gorm:""`                     // 内容
	CreatedAt time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time      `gorm:""`                     // 更新日時
	DeletedAt gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Threads []*Thread

type NewThreadParams struct {
	UserType  int32
	ContactID string
	Content   string
}

func NewThread(params *NewThreadParams) *Thread {
	return &Thread{
		ID:        uuid.Base58Encode(uuid.New()),
		ContactID: params.ContactID,
		UserType:  params.UserType,
		Content:   params.Content,
	}
}

func (t *Thread) Fill(userID string) {
	if userID != "" {
		t.UserID = userID
	}
}
