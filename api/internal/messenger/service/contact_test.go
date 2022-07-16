package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListContacts(t *testing.T) {
	t.Parallel()

	now := jst.Date(20222, 7, 7, 18, 30, 0, 0)
	params := &database.ListContactsParams{
		Limit:  20,
		Offset: 0,
	}
	contacts := entity.Contacts{
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
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *messenger.ListContactsInput
		expect      entity.Contacts
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().List(gomock.Any(), params).Return(contacts, nil)
				mocks.db.Contact.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &messenger.ListContactsInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      contacts,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.ListContactsInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list contacts",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().List(gomock.Any(), params).Return(nil, errmock)
				mocks.db.Contact.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &messenger.ListContactsInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
		{
			name: "failed to count contacts",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().List(gomock.Any(), params).Return(contacts, nil)
				mocks.db.Contact.EXPECT().Count(gomock.Any(), params).Return(int64(0), errmock)
			},
			input: &messenger.ListContactsInput{
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListContacts(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestGetContact(t *testing.T) {
	t.Parallel()

	now := jst.Date(20222, 7, 7, 18, 30, 0, 0)
	contact := &entity.Contact{
		ID:          "contact-id",
		Title:       "お問い合わせ件名",
		Content:     "お問い合わせ内容です。",
		Username:    "あんど どっと",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+819012345678",
		Status:      entity.ContactStatusInprogress,
		Priority:    entity.ContactPriorityMiddle,
		Note:        "対応者のメモです",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.GetContactInput
		expect    *entity.Contact
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Get(ctx, "contact-id").Return(contact, nil)
			},
			input: &messenger.GetContactInput{
				ContactID: "contact-id",
			},
			expect:    contact,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.GetContactInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get contact",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Get(ctx, "contact-id").Return(nil, errmock)
			},
			input: &messenger.GetContactInput{
				ContactID: "contact-id",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
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
	contact := &entity.Contact{
		ID:          "contact-id",
		Title:       "お問い合わせ件名",
		Content:     "お問い合わせ内容",
		Username:    "お問い合わせ氏名",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+819012345678",
		Status:      entity.ContactStatusUnknown,
		Priority:    entity.ContactPriorityUnknown,
		CreatedAt:   jst.Date(2022, 7, 13, 18, 30, 0, 0),
		UpdatedAt:   jst.Date(2022, 7, 13, 18, 30, 0, 0),
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.CreateContactInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, contact *entity.Contact) error {
						expect := &entity.Contact{
							ID:          contact.ID, // ignore
							Title:       "お問い合わせ件名",
							Content:     "お問い合わせ内容",
							Username:    "お問い合わせ氏名",
							Email:       "test-user@and-period.jp",
							PhoneNumber: "+819012345678",
							Status:      entity.ContactStatusUnknown,
							Priority:    entity.ContactPriorityUnknown,
						}
						assert.Equal(t, expect, contact)
						return nil
					})
				mocks.db.Contact.EXPECT().Get(gomock.Any(), gomock.Any()).Return(contact, nil)
				mocks.db.ReceivedQueue.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return("", nil)
			},
			input: &messenger.CreateContactInput{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
			},
			expectErr: nil,
		},
		{
			name: "success without notify",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.db.Contact.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errmock)
			},
			input: &messenger.CreateContactInput{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.CreateContactInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create contact",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Create(ctx, gomock.Any()).Return(errmock)
			},
			input: &messenger.CreateContactInput{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
			},
			expectErr: exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateContact(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateContact(t *testing.T) {
	t.Parallel()

	params := &database.UpdateContactParams{
		Status:   entity.ContactStatusInprogress,
		Priority: entity.ContactPriorityMiddle,
		Note:     "テストです。",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.UpdateContactInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Update(ctx, "contact-id", params).Return(nil)
			},
			input: &messenger.UpdateContactInput{
				ContactID: "contact-id",
				Status:    entity.ContactStatusInprogress,
				Priority:  entity.ContactPriorityMiddle,
				Note:      "テストです。",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.UpdateContactInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update contact",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Contact.EXPECT().Update(ctx, "contact-id", params).Return(errmock)
			},
			input: &messenger.UpdateContactInput{
				ContactID: "contact-id",
				Status:    entity.ContactStatusInprogress,
				Priority:  entity.ContactPriorityMiddle,
				Note:      "テストです。",
			},
			expectErr: exception.ErrUnknown,
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
