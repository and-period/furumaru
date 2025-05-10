package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceivedQueues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload *WorkerPayload
		expect  ReceivedQueues
	}{
		{
			name: "success",
			payload: &WorkerPayload{
				QueueID:   "id",
				EventType: EventTypeRegisterAdmin,
				UserType:  UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Email: &MailConfig{
					TemplateID:    EmailTemplateIDAdminRegister,
					Substitutions: map[string]interface{}{"パスワード": "!Qaz2wsx"},
				},
				Message: &MessageConfig{},
				Push:    &PushConfig{},
				Report:  &ReportConfig{},
			},
			expect: ReceivedQueues{
				{
					ID:         "id",
					NotifyType: NotifyTypeEmail,
					EventType:  EventTypeRegisterAdmin,
					UserType:   UserTypeAdmin,
					UserIDs:    []string{"admin-id"},
					Done:       false,
				},
				{
					ID:         "id",
					NotifyType: NotifyTypeMessage,
					EventType:  EventTypeRegisterAdmin,
					UserType:   UserTypeAdmin,
					UserIDs:    []string{"admin-id"},
					Done:       false,
				},
				{
					ID:         "id",
					NotifyType: NotifyTypePush,
					EventType:  EventTypeRegisterAdmin,
					UserType:   UserTypeAdmin,
					UserIDs:    []string{"admin-id"},
					Done:       false,
				},
				{
					ID:         "id",
					NotifyType: NotifyTypeReport,
					EventType:  EventTypeRegisterAdmin,
					UserType:   UserTypeAdmin,
					UserIDs:    []string{"admin-id"},
					Done:       false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewReceivedQueues(tt.payload)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
