package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// MessageType - メッセージ種別
type MessageType int32

const (
	MessageTypeUnknown      MessageType = 0
	MessageTypeNotification MessageType = 1 // お知らせ
)

// Message - メッセージ情報
type Message struct {
	ID         string      `gorm:"primaryKey;<-:create"` // メッセージID
	UserType   UserType    `gorm:""`                     // ユーザー種別
	UserID     string      `gorm:""`                     // ユーザーID
	Type       MessageType `gorm:""`                     // メッセージ種別
	Title      string      `gorm:""`                     // メッセージ件名
	Body       string      `gorm:""`                     // メッセージ内容
	Link       string      `gorm:""`                     // 遷移先リンク
	Read       bool        `gorm:""`                     // 既読フラグ
	ReceivedAt time.Time   `gorm:"<-:create"`            // 受信日時
	CreatedAt  time.Time   `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time   `gorm:""`                     // 更新日時
}

type Messages []*Message

type NewMessageParams struct {
	UserType   UserType
	UserID     string
	Type       MessageType
	Title      string
	Body       string
	Link       string
	ReceivedAt time.Time
}

func NewMessage(params *NewMessageParams) *Message {
	var userType UserType
	switch params.UserType {
	case UserTypeUser:
		userType = UserTypeUser
	case UserTypeAdmin,
		UserTypeAdministrator,
		UserTypeCoordinator,
		UserTypeProducer:
		userType = UserTypeAdmin
	}
	return &Message{
		ID:         uuid.Base58Encode(uuid.New()),
		UserType:   userType,
		UserID:     params.UserID,
		Type:       params.Type,
		Title:      params.Title,
		Body:       params.Body,
		Link:       params.Link,
		Read:       false,
		ReceivedAt: params.ReceivedAt,
	}
}

type NewMessagesParams struct {
	UserType   UserType
	UserIDs    []string
	Type       MessageType
	Title      string
	Body       string
	Link       string
	ReceivedAt time.Time
}

func NewMessages(params *NewMessagesParams) Messages {
	res := make(Messages, len(params.UserIDs))
	for i := range params.UserIDs {
		p := &NewMessageParams{
			UserType:   params.UserType,
			UserID:     params.UserIDs[i],
			Type:       params.Type,
			Title:      params.Title,
			Body:       params.Body,
			Link:       params.Link,
			ReceivedAt: params.ReceivedAt,
		}
		res[i] = NewMessage(p)
	}
	return res
}
