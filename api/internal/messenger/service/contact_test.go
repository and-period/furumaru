package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/stretchr/testify/assert"
)

func TestListContacts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.ListContactsInput
		expect    entity.Contacts
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &messenger.ListContactsInput{
				Limit:  20,
				Offset: 0,
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.ListContactsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListContacts(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetContact(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.GetContactInput
		expect    *entity.Contact
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &messenger.GetContactInput{
				ContactID: "contact-id",
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.GetContactInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetContact(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateContact(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.CreateContactInput
		expect    *entity.Contact
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &messenger.CreateContactInput{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.CreateContactInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CreateContact(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUpdateContact(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.UpdateContactInput
		expect    *entity.Contact
		expectErr error
	}{
		{
			name:  "not implemented",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &messenger.UpdateContactInput{
				ContactID: "contact-id",
				Status:    entity.ContactStatusInprogress,
				Priority:  entity.ContactPriorityMiddle,
				Note:      "テストです。",
			},
			expect:    nil,
			expectErr: exception.ErrNotImplemented,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.UpdateContactInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateContact(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
