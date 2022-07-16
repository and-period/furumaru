package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/messenger"
	mentity "github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestListContacts(t *testing.T) {
	t.Parallel()
	contactsIn := &messenger.ListContactsInput{
		Limit:  20,
		Offset: 0,
	}
	contacts := mentity.Contacts{
		{
			ID:          "contact-id",
			Title:       "お問い合わせ件名",
			Content:     "お問い合わせ内容です。",
			Username:    "あんど どっと",
			Email:       "test-user@and-period.jp",
			PhoneNumber: "+819012345678",
			Status:      mentity.ContactStatusInprogress,
			Priority:    mentity.ContactPriorityMiddle,
			Note:        "対応者のメモです",
			CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		query  string
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().ListContacts(gomock.Any(), contactsIn).Return(contacts, int64(1), nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ContactsResponse{
					Contacts: []*response.Contact{
						{
							ID:          "contact-id",
							Title:       "お問い合わせ件名",
							Content:     "お問い合わせ内容です。",
							Username:    "あんど どっと",
							Email:       "test-user@and-period.jp",
							PhoneNumber: "+819012345678",
							Status:      int32(service.ContactStatusInprogress),
							Priority:    int32(service.ContactPriorityMiddle),
							Note:        "対応者のメモです",
							CreatedAt:   1640962800,
							UpdatedAt:   1640962800,
						},
					},
					Total: 1,
				},
			},
		},
		{
			name:  "invalid limit",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid offset",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to list contacts",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().ListContacts(gomock.Any(), contactsIn).Return(nil, int64(0), errmock)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const prefix = "/v1/contacts"
			path := fmt.Sprintf("%s%s", prefix, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetContact(t *testing.T) {
	t.Parallel()
	contactIn := &messenger.GetContactInput{
		ContactID: "contact-id",
	}
	contact := &mentity.Contact{
		ID:          "contact-id",
		Title:       "お問い合わせ件名",
		Content:     "お問い合わせ内容です。",
		Username:    "あんど どっと",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+819012345678",
		Status:      mentity.ContactStatusInprogress,
		Priority:    mentity.ContactPriorityMiddle,
		Note:        "対応者のメモです",
		CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		contactID string
		expect    *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().GetContact(gomock.Any(), contactIn).Return(contact, nil)
			},
			contactID: "contact-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ContactResponse{
					Contact: &response.Contact{
						ID:          "contact-id",
						Title:       "お問い合わせ件名",
						Content:     "お問い合わせ内容です。",
						Username:    "あんど どっと",
						Email:       "test-user@and-period.jp",
						PhoneNumber: "+819012345678",
						Status:      int32(service.ContactStatusInprogress),
						Priority:    int32(service.ContactPriorityMiddle),
						Note:        "対応者のメモです",
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
		},
		{
			name: "failed to get contact",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().GetContact(gomock.Any(), contactIn).Return(nil, errmock)
			},
			contactID: "contact-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const prefix = "/v1/contacts"
			path := fmt.Sprintf("%s/%s", prefix, tt.contactID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestUpdateContact(t *testing.T) {
	t.Parallel()
	contactIn := &messenger.UpdateContactInput{
		ContactID: "contact-id",
		Status:    mentity.ContactStatusInprogress,
		Priority:  mentity.ContactPriorityMiddle,
		Note:      "対応者のメモです",
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		contactID string
		req       *request.UpdateContactRequest
		expect    *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().UpdateContact(gomock.Any(), contactIn).Return(nil)
			},
			contactID: "contact-id",
			req: &request.UpdateContactRequest{
				Status:   int32(service.ContactStatusInprogress),
				Priority: int32(service.ContactPriorityMiddle),
				Note:     "対応者のメモです",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update contact",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().UpdateContact(gomock.Any(), contactIn).Return(errmock)
			},
			contactID: "contact-id",
			req: &request.UpdateContactRequest{
				Status:   int32(service.ContactStatusInprogress),
				Priority: int32(service.ContactPriorityMiddle),
				Note:     "対応者のメモです",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const prefix = "/v1/contacts"
			path := fmt.Sprintf("%s/%s", prefix, tt.contactID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}
