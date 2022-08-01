package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestListMessages(t *testing.T) {
	t.Parallel()

	in := &messenger.ListMessagesInput{
		UserType: entity.UserTypeAdmin,
		UserID:   idmock,
		Limit:    20,
		Offset:   0,
	}
	messages := entity.Messages{
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
				mocks.messenger.EXPECT().ListMessages(gomock.Any(), in).Return(messages, int64(1), nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.MessagesResponse{
					Messages: []*response.Message{
						{
							ID:         "message-id",
							Type:       int32(service.MessageTypeNotification),
							Title:      "メッセージタイトル",
							Body:       "メッセージの内容です。",
							Link:       "https://and-period.jp",
							Read:       false,
							ReceivedAt: 1640962800,
							CreatedAt:  1640962800,
							UpdatedAt:  1640962800,
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
			name: "failed to list messages",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().ListMessages(gomock.Any(), in).Return(nil, int64(0), errmock)
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
			const format = "/v1/messages%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetMessage(t *testing.T) {
	t.Parallel()

	in := &messenger.GetMessageInput{
		MessageID: "message-id",
		UserType:  entity.UserTypeAdmin,
		UserID:    idmock,
	}
	message := &entity.Message{
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
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		messageID string
		expect    *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().GetMessage(gomock.Any(), in).Return(message, nil)
			},
			messageID: "message-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.MessageResponse{
					Message: &response.Message{
						ID:         "message-id",
						Type:       int32(service.MessageTypeNotification),
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
		{
			name: "failed to get message",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().GetMessage(gomock.Any(), in).Return(nil, errmock)
			},
			messageID: "message-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/messages/%s"
			path := fmt.Sprintf(format, tt.messageID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}
