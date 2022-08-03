package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestMessageConfig_Fields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		message *MessageConfig
		expect  map[string]string
	}{
		{
			name: "success",
			message: &MessageConfig{
				MessageID:   MessageIDNotification,
				MessageType: MessageTypeNotification,
				Title:       "メッセージのタイトル",
				Author:      "&.スタッフ",
				Link:        "https://and-period.jp",
				ReceivedAt:  jst.Date(2022, 7, 14, 18, 30, 0, 0),
			},
			expect: map[string]string{
				"Title":  "メッセージのタイトル",
				"Author": "&.スタッフ",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.message.Fields())
		})
	}
}
