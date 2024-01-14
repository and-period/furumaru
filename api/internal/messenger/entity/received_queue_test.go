package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestReceivedQueue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload *WorkerPayload
		expect  *ReceivedQueue
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
			},
			expect: &ReceivedQueue{
				ID:        "id",
				EventType: EventTypeRegisterAdmin,
				UserType:  UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Done:      false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewReceivedQueue(tt.payload)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestReceivedQueue_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		queue  *ReceivedQueue
		expect *ReceivedQueue
		hasErr bool
	}{
		{
			name: "success",
			queue: &ReceivedQueue{
				ID:          "id",
				EventType:   EventTypeRegisterAdmin,
				UserType:    UserTypeAdmin,
				UserIDsJSON: datatypes.JSON([]byte(`["admin-id"]`)),
				Done:        false,
				CreatedAt:   jst.Date(2022, 7, 10, 18, 30, 0, 0),
				UpdatedAt:   jst.Date(2022, 7, 10, 18, 30, 0, 0),
			},
			expect: &ReceivedQueue{
				ID:          "id",
				EventType:   EventTypeRegisterAdmin,
				UserType:    UserTypeAdmin,
				UserIDs:     []string{"admin-id"},
				UserIDsJSON: datatypes.JSON([]byte(`["admin-id"]`)),
				Done:        false,
				CreatedAt:   jst.Date(2022, 7, 10, 18, 30, 0, 0),
				UpdatedAt:   jst.Date(2022, 7, 10, 18, 30, 0, 0),
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.queue.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.queue)
		})
	}
}

func TestReceivedQueue_FillJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		queue  *ReceivedQueue
		expect *ReceivedQueue
		hasErr bool
	}{
		{
			name: "success",
			queue: &ReceivedQueue{
				ID:        "id",
				EventType: EventTypeRegisterAdmin,
				UserType:  UserTypeAdmin,
				UserIDs:   []string{"admin-id"},
				Done:      false,
				CreatedAt: jst.Date(2022, 7, 10, 18, 30, 0, 0),
				UpdatedAt: jst.Date(2022, 7, 10, 18, 30, 0, 0),
			},
			expect: &ReceivedQueue{
				ID:          "id",
				EventType:   EventTypeRegisterAdmin,
				UserType:    UserTypeAdmin,
				UserIDs:     []string{"admin-id"},
				UserIDsJSON: datatypes.JSON([]byte(`["admin-id"]`)),
				Done:        false,
				CreatedAt:   jst.Date(2022, 7, 10, 18, 30, 0, 0),
				UpdatedAt:   jst.Date(2022, 7, 10, 18, 30, 0, 0),
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.queue.FillJSON()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.queue)
		})
	}
}
