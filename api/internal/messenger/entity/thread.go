package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

type ThreadUserType int32

const (
	ThreadUserTypeUnknown ThreadUserType = iota // 不明
	ThreadUserTypeAdmin                         // 管理者
	ThreadUserTypeUser                          // ユーザー
	ThreadUserTypeGuest                         // ゲスト(ユーザIDなし)
)

// Thread - お問い合わせ会話履歴
type Thread struct {
	ID        string         `gorm:"primaryKey;<-:create"` // お問い合わせ会話履歴ID
	ContactID string         `gorm:""`                     // お問い合わせID
	UserID    string         `gorm:"default:null"`         // 送信者ID(ゲストの場合null)
	UserType  ThreadUserType `gorm:""`                     // 送信者の種別(不明:0, admin:1, user:2, guest:3)
	Content   string         `gorm:""`                     // 内容
	CreatedAt time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time      `gorm:""`                     // 更新日時
	DeletedAt gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Threads []*Thread

type NewThreadParams struct {
	UserType  ThreadUserType
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

func (ts Threads) IDs() []string {
	return set.UniqBy(ts, func(t *Thread) string {
		return t.ID
	})
}

func (ts Threads) UserIDs() []string {
	return set.UniqBy(ts, func(t *Thread) string {
		if t.UserType != ThreadUserTypeUser {
			return ""
		}
		return t.UserID
	})
}

func (ts Threads) AdminIDs() []string {
	return set.UniqBy(ts, func(t *Thread) string {
		if t.UserType != ThreadUserTypeAdmin {
			return ""
		}
		return t.UserID
	})
}

func (ts Threads) Fill() {
	userIDs := ts.UserIDs()
	for _, t := range ts {
		for _, userID := range userIDs {
			if t.UserID == userID {
				t.UserID = userID
			}
		}
	}
}
