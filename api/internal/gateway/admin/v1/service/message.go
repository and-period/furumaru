package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

// MessageType - メッセージ種別
type MessageType int32

const (
	MessageTypeUnknown      MessageType = 0
	MessageTypeNotification MessageType = 1 // お知らせ
)

type Message struct {
	types.Message
}

type Messages []*Message

func NewMessageType(typ entity.MessageType) MessageType {
	switch typ {
	case entity.MessageTypeNotification:
		return MessageTypeNotification
	default:
		return MessageTypeUnknown
	}
}

func (t MessageType) Response() int32 {
	return int32(t)
}

func NewMessage(message *entity.Message) *Message {
	return &Message{
		Message: types.Message{
			ID:         message.ID,
			Type:       NewMessageType(message.Type).Response(),
			Title:      message.Title,
			Body:       message.Body,
			Link:       message.Link,
			Read:       message.Read,
			ReceivedAt: message.ReceivedAt.Unix(),
			CreatedAt:  message.CreatedAt.Unix(),
			UpdatedAt:  message.UpdatedAt.Unix(),
		},
	}
}

func (m *Message) Response() *types.Message {
	return &m.Message
}

func NewMessages(messages entity.Messages) Messages {
	res := make(Messages, len(messages))
	for i := range messages {
		res[i] = NewMessage(messages[i])
	}
	return res
}

func (ms Messages) Response() []*types.Message {
	res := make([]*types.Message, len(ms))
	for i := range ms {
		res[i] = ms[i].Response()
	}
	return res
}
