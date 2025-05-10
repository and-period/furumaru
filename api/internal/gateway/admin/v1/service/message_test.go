package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestMessageType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		messageType entity.MessageType
		expect      MessageType
	}{
		{
			name:        "notification",
			messageType: entity.MessageTypeNotification,
			expect:      MessageTypeNotification,
		},
		{
			name:        "unknown",
			messageType: entity.MessageTypeUnknown,
			expect:      MessageTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewMessageType(tt.messageType))
		})
	}
}

func TestMessageType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		messageType MessageType
		expect      int32
	}{
		{
			name:        "success",
			messageType: MessageTypeNotification,
			expect:      1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.messageType.Response())
		})
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		message *entity.Message
		expect  *Message
	}{
		{
			name: "success",
			message: &entity.Message{
				ID:         "message-id",
				UserType:   entity.UserTypeUser,
				UserID:     "user-id",
				Type:       entity.MessageTypeNotification,
				Title:      "メッセージタイトル",
				Body:       "メッセージの内容です。",
				Link:       "https://and-period.jp",
				Read:       false,
				ReceivedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Message{
				Message: response.Message{
					ID:         "message-id",
					Type:       int32(MessageTypeNotification),
					Title:      "メッセージタイトル",
					Body:       "メッセージの内容です。",
					Link:       "https://and-period.jp",
					Read:       false,
					ReceivedAt: 1640962800,
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewMessage(tt.message))
		})
	}
}

func TestMessage_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		message *Message
		expect  *response.Message
	}{
		{
			name: "success",
			message: &Message{
				Message: response.Message{
					ID:         "message-id",
					Type:       int32(MessageTypeNotification),
					Title:      "メッセージタイトル",
					Body:       "メッセージの内容です。",
					Link:       "https://and-period.jp",
					Read:       false,
					ReceivedAt: 1640962800,
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
			expect: &response.Message{
				ID:         "message-id",
				Type:       int32(MessageTypeNotification),
				Title:      "メッセージタイトル",
				Body:       "メッセージの内容です。",
				Link:       "https://and-period.jp",
				Read:       false,
				ReceivedAt: 1640962800,
				CreatedAt:  1640962800,
				UpdatedAt:  1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.message.Response())
		})
	}
}

func TestMessages(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		messages entity.Messages
		expect   Messages
	}{
		{
			name: "success",
			messages: entity.Messages{
				{
					ID:         "message-id",
					UserType:   entity.UserTypeUser,
					UserID:     "user-id",
					Type:       entity.MessageTypeNotification,
					Title:      "メッセージタイトル",
					Body:       "メッセージの内容です。",
					Link:       "https://and-period.jp",
					Read:       false,
					ReceivedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Messages{
				{
					Message: response.Message{
						ID:         "message-id",
						Type:       int32(MessageTypeNotification),
						Title:      "メッセージタイトル",
						Body:       "メッセージの内容です。",
						Link:       "https://and-period.jp",
						Read:       false,
						ReceivedAt: 1640962800,
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewMessages(tt.messages))
		})
	}
}

func TestMessages_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		messages Messages
		expect   []*response.Message
	}{
		{
			name: "success",
			messages: Messages{
				{
					Message: response.Message{
						ID:         "message-id",
						Type:       int32(MessageTypeNotification),
						Title:      "メッセージタイトル",
						Body:       "メッセージの内容です。",
						Link:       "https://and-period.jp",
						Read:       false,
						ReceivedAt: 1640962800,
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
					},
				},
			},
			expect: []*response.Message{
				{
					ID:         "message-id",
					Type:       int32(MessageTypeNotification),
					Title:      "メッセージタイトル",
					Body:       "メッセージの内容です。",
					Link:       "https://and-period.jp",
					Read:       false,
					ReceivedAt: 1640962800,
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.messages.Response())
		})
	}
}
