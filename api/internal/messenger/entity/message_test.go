package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		params *NewMessageParams
		expect *Message
	}{
		{
			name: "success user message",
			params: &NewMessageParams{
				UserType:   UserTypeUser,
				UserID:     "user-id",
				Type:       MessageTypeNotification,
				Title:      "メッセージタイトル",
				Body:       "メッセージの内容です。",
				Link:       "https://and-period.jp",
				ReceivedAt: now,
			},
			expect: &Message{
				UserType:   UserTypeUser,
				UserID:     "user-id",
				Type:       MessageTypeNotification,
				Title:      "メッセージタイトル",
				Body:       "メッセージの内容です。",
				Link:       "https://and-period.jp",
				Read:       false,
				ReceivedAt: now,
			},
		},
		{
			name: "success admin message",
			params: &NewMessageParams{
				UserType:   UserTypeAdmin,
				UserID:     "user-id",
				Type:       MessageTypeNotification,
				Title:      "メッセージタイトル",
				Body:       "メッセージの内容です。",
				Link:       "https://and-period.jp",
				ReceivedAt: now,
			},
			expect: &Message{
				UserType:   UserTypeAdmin,
				UserID:     "user-id",
				Type:       MessageTypeNotification,
				Title:      "メッセージタイトル",
				Body:       "メッセージの内容です。",
				Link:       "https://and-period.jp",
				Read:       false,
				ReceivedAt: now,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewMessage(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestMessages(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		params *NewMessagesParams
		expect Messages
	}{
		{
			name: "success",
			params: &NewMessagesParams{
				UserType:   UserTypeUser,
				UserIDs:    []string{"user-id"},
				Type:       MessageTypeNotification,
				Title:      "メッセージタイトル",
				Body:       "メッセージの内容です。",
				Link:       "https://and-period.jp",
				ReceivedAt: now,
			},
			expect: Messages{
				{
					UserType:   UserTypeUser,
					UserID:     "user-id",
					Type:       MessageTypeNotification,
					Title:      "メッセージタイトル",
					Body:       "メッセージの内容です。",
					Link:       "https://and-period.jp",
					Read:       false,
					ReceivedAt: now,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewMessages(tt.params)
			for i := range actual {
				actual[i].ID = "" // ignore
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}
