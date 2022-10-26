package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestListUsers(t *testing.T) {
	t.Parallel()

	usersIn := &user.ListUsersInput{
		Limit:  20,
		Offset: 0,
	}
	users := uentity.Users{
		{
			ID:         "user-id",
			Registered: true,
			CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			Member: entity.Member{
				AccountID:    "",
				CognitoID:    "cognito-id",
				Username:     "",
				ProviderType: entity.ProviderTypeEmail,
				Email:        "test-user@and-period.jp",
				PhoneNumber:  "+810000000000",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
				CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				VerifiedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
		},
	}
	ordersIn := &store.AggregateOrdersInput{
		UserIDs: []string{"user-id"},
	}
	orders := sentity.AggregatedOrders{
		{
			UserID:     "user-id",
			OrderCount: 2,
			Subtotal:   6000,
			Discount:   1000,
		},
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		query   string
		expect  *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListUsers(gomock.Any(), usersIn).Return(users, int64(1), nil)
				mocks.store.EXPECT().AggregateOrders(gomock.Any(), ordersIn).Return(orders, nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.UsersResponse{
					Users: []*response.UserList{
						{
							ID:          "user-id",
							Lastname:    "",
							Firstname:   "",
							Registered:  true,
							Address:     "",
							TotalOrder:  2,
							TotalAmount: 6000,
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/users%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path, tt.options...)
		})
	}
}
