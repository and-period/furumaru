package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestContactStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ContactStatus
		expect ContactStatus
	}{
		{
			name:   "todo",
			status: entity.ContactStatusToDo,
			expect: ContactStatusToDo,
		},
		{
			name:   "inprogress",
			status: entity.ContactStatusInprogress,
			expect: ContactStatusInprogress,
		},
		{
			name:   "done",
			status: entity.ContactStatusDone,
			expect: ContactStatusDone,
		},
		{
			name:   "discard",
			status: entity.ContactStatusDiscard,
			expect: ContactStatusDiscard,
		},
		{
			name:   "unknown",
			status: entity.ContactStatusUnknown,
			expect: ContactStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewContactStatus(tt.status))
		})
	}
}

func TestContactStatus_MessengerEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ContactStatus
		expect entity.ContactStatus
	}{
		{
			name:   "todo",
			status: ContactStatusToDo,
			expect: entity.ContactStatusToDo,
		},
		{
			name:   "inprogress",
			status: ContactStatusInprogress,
			expect: entity.ContactStatusInprogress,
		},
		{
			name:   "done",
			status: ContactStatusDone,
			expect: entity.ContactStatusDone,
		},
		{
			name:   "discard",
			status: ContactStatusDiscard,
			expect: entity.ContactStatusDiscard,
		},
		{
			name:   "unknown",
			status: ContactStatusUnknown,
			expect: entity.ContactStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.MessengerEntity())
		})
	}
}

func TestContactStatus_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ContactStatus
		expect int32
	}{
		{
			name:   "success",
			status: ContactStatusInprogress,
			expect: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestContactPriority(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		priority entity.ContactPriority
		expect   ContactPriority
	}{
		{
			name:     "low",
			priority: entity.ContactPriorityLow,
			expect:   ContactPriorityLow,
		},
		{
			name:     "middle",
			priority: entity.ContactPriorityMiddle,
			expect:   ContactPriorityMiddle,
		},
		{
			name:     "high",
			priority: entity.ContactPriorityHigh,
			expect:   ContactPriorityHigh,
		},
		{
			name:     "unknown",
			priority: entity.ContactPriorityUnknown,
			expect:   ContactPriorityUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewContactPriority(tt.priority))
		})
	}
}

func TestContactPriority_MessengerEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		priority ContactPriority
		expect   entity.ContactPriority
	}{
		{
			name:     "low",
			priority: ContactPriorityLow,
			expect:   entity.ContactPriorityLow,
		},
		{
			name:     "middle",
			priority: ContactPriorityMiddle,
			expect:   entity.ContactPriorityMiddle,
		},
		{
			name:     "high",
			priority: ContactPriorityHigh,
			expect:   entity.ContactPriorityHigh,
		},
		{
			name:     "unknown",
			priority: ContactPriorityUnknown,
			expect:   entity.ContactPriorityUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.priority.MessengerEntity())
		})
	}
}

func TestContactPriority_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		priority ContactPriority
		expect   int32
	}{
		{
			name:     "success",
			priority: ContactPriorityMiddle,
			expect:   2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.priority.Response())
		})
	}
}

func TestContact(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		contact *entity.Contact
		expect  *Contact
	}{
		{
			name: "success",
			contact: &entity.Contact{
				ID:          "contact-id",
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容です。",
				Username:    "あんど どっと",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      entity.ContactStatusInprogress,
				Priority:    entity.ContactPriorityMiddle,
				Note:        "対応者のメモです",
				CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Contact{
				Contact: response.Contact{
					ID:          "contact-id",
					Title:       "お問い合わせ件名",
					Content:     "お問い合わせ内容です。",
					Username:    "あんど どっと",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					Status:      int32(ContactStatusInprogress),
					Priority:    int32(ContactPriorityMiddle),
					Note:        "対応者のメモです",
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewContact(tt.contact))
		})
	}
}

func TestContact_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		contact *Contact
		expect  *response.Contact
	}{
		{
			name: "success",
			contact: &Contact{
				Contact: response.Contact{
					ID:          "contact-id",
					Title:       "お問い合わせ件名",
					Content:     "お問い合わせ内容です。",
					Username:    "あんど どっと",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					Status:      int32(ContactStatusInprogress),
					Priority:    int32(ContactPriorityMiddle),
					Note:        "対応者のメモです",
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
			expect: &response.Contact{
				ID:          "contact-id",
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容です。",
				Username:    "あんど どっと",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
				Status:      int32(ContactStatusInprogress),
				Priority:    int32(ContactPriorityMiddle),
				Note:        "対応者のメモです",
				CreatedAt:   1640962800,
				UpdatedAt:   1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.contact.Response())
		})
	}
}

func TestContacts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts entity.Contacts
		expect   Contacts
	}{
		{
			name: "success",
			contacts: entity.Contacts{
				{
					ID:          "contact-id",
					Title:       "お問い合わせ件名",
					Content:     "お問い合わせ内容です。",
					Username:    "あんど どっと",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					Status:      entity.ContactStatusInprogress,
					Priority:    entity.ContactPriorityMiddle,
					Note:        "対応者のメモです",
					CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Contacts{
				{
					Contact: response.Contact{
						ID:          "contact-id",
						Title:       "お問い合わせ件名",
						Content:     "お問い合わせ内容です。",
						Username:    "あんど どっと",
						Email:       "test-user@and-period.jp",
						PhoneNumber: "+819012345678",
						Status:      int32(ContactStatusInprogress),
						Priority:    int32(ContactPriorityMiddle),
						Note:        "対応者のメモです",
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewContacts(tt.contacts))
		})
	}
}

func TestContacts_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contacts Contacts
		expect   []*response.Contact
	}{
		{
			name: "success",
			contacts: Contacts{
				{
					Contact: response.Contact{
						ID:          "contact-id",
						Title:       "お問い合わせ件名",
						Content:     "お問い合わせ内容です。",
						Username:    "あんど どっと",
						Email:       "test-user@and-period.jp",
						PhoneNumber: "+819012345678",
						Status:      int32(ContactStatusInprogress),
						Priority:    int32(ContactPriorityMiddle),
						Note:        "対応者のメモです",
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
			expect: []*response.Contact{
				{
					ID:          "contact-id",
					Title:       "お問い合わせ件名",
					Content:     "お問い合わせ内容です。",
					Username:    "あんど どっと",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					Status:      int32(ContactStatusInprogress),
					Priority:    int32(ContactPriorityMiddle),
					Note:        "対応者のメモです",
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.contacts.Response())
		})
	}
}
