package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/golang/mock/gomock"
)

func TestCreateContact(t *testing.T) {
	t.Parallel()
	in := &messenger.CreateContactInput{
		Title:       "お問い合わせ件名",
		Content:     "お問い合わせ内容",
		Username:    "お問い合わせ氏名",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+819012345678",
	}
	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateContactRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().CreateContact(gomock.Any(), in).Return(nil, nil)
			},
			req: &request.CreateContactRequest{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to create contact",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().CreateContact(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateContactRequest{
				Title:       "お問い合わせ件名",
				Content:     "お問い合わせ内容",
				Username:    "お問い合わせ氏名",
				Email:       "test-user@and-period.jp",
				PhoneNumber: "+819012345678",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/contacts"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}
