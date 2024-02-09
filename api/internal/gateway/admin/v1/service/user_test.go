package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestUserStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		status   entity.UserStatus
		expect   UserStatus
		response int32
	}{
		{
			name:     "guest",
			status:   entity.UserStatusGuest,
			expect:   UserStatusGuest,
			response: 1,
		},
		{
			name:     "provisional",
			status:   entity.UserStatusProvisional,
			expect:   UserStatusProvisional,
			response: 2,
		},
		{
			name:     "verified",
			status:   entity.UserStatusVerified,
			expect:   UserStatusVerified,
			response: 3,
		},
		{
			name:     "deactivated",
			status:   entity.UserStatusDeactivated,
			expect:   UserStatusDeactivated,
			response: 4,
		},
		{
			name:     "unknown",
			status:   entity.UserStatusUnknown,
			expect:   UserStatusUnknown,
			response: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewUserStatus(tt.status)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.response, actual.Response())
		})
	}
}

func TestUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		user    *uentity.User
		address *uentity.Address
		expect  *User
	}{
		{
			name: "success member",
			user: &entity.User{
				ID:         "user-id",
				Registered: true,
				Status:     entity.UserStatusVerified,
				Member: uentity.Member{
					UserID:        "user-id",
					AccountID:     "account-id",
					CognitoID:     "cognito-id",
					Username:      "username",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
					ProviderType:  uentity.ProviderTypeEmail,
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					VerifiedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			address: &entity.Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: uentity.AddressRevision{
					ID:             1,
					Lastname:       "&.",
					Firstname:      "購入者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "こうにゅうしゃ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
				},
			},
			expect: &User{
				User: response.User{
					ID:            "user-id",
					Status:        int32(UserStatusVerified),
					Registered:    true,
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				address: Address{
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					revisionID: 1,
				},
			},
		},
		{
			name: "success guest",
			user: &entity.User{
				ID:         "user-id",
				Registered: false,
				Status:     entity.UserStatusGuest,
				Guest: uentity.Guest{
					UserID:      "user-id",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			address: &entity.Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: uentity.AddressRevision{
					ID:             1,
					Lastname:       "&.",
					Firstname:      "購入者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "こうにゅうしゃ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
				},
			},
			expect: &User{
				User: response.User{
					ID:            "user-id",
					Status:        int32(UserStatusGuest),
					Registered:    false,
					Username:      "",
					AccountID:     "",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				address: Address{
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					revisionID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewUser(tt.user, tt.address))
		})
	}
}

func TestUser_Address(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *User
		expect *Address
	}{
		{
			name: "success",
			user: &User{
				User: response.User{
					ID:            "user-id",
					Status:        int32(UserStatusVerified),
					Registered:    true,
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				address: Address{
					Address: response.Address{
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					revisionID: 1,
				},
			},
			expect: &Address{
				Address: response.Address{
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "こうにゅうしゃ",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
				},
				revisionID: 1,
			},
		},
		{
			name:   "empty",
			user:   &User{},
			expect: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Address())
		})
	}
}

func TestUser_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *User
		expect *response.User
	}{
		{
			name: "success",
			user: &User{
				User: response.User{
					ID:            "user-id",
					Status:        int32(UserStatusVerified),
					Registered:    true,
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				address: Address{
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					revisionID: 1,
				},
			},
			expect: &response.User{
				ID:            "user-id",
				Status:        int32(UserStatusVerified),
				Registered:    true,
				Username:      "username",
				AccountID:     "account-id",
				Lastname:      "&.",
				Firstname:     "購入者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "こうにゅうしゃ",
				Email:         "test-user@and-period.jp",
				PhoneNumber:   "+819012345678",
				CreatedAt:     1640962800,
				UpdatedAt:     1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Response())
		})
	}
}

func TestUsers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		users     uentity.Users
		addresses map[string]*uentity.Address
		expect    Users
	}{
		{
			name: "success",
			users: uentity.Users{
				{
					ID:         "user-id",
					Registered: true,
					Status:     entity.UserStatusVerified,
					Member: uentity.Member{
						UserID:        "user-id",
						AccountID:     "account-id",
						CognitoID:     "cognito-id",
						Username:      "username",
						Lastname:      "&.",
						Firstname:     "購入者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "こうにゅうしゃ",
						ProviderType:  uentity.ProviderTypeEmail,
						Email:         "test-user@and-period.jp",
						PhoneNumber:   "+819012345678",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
						VerifiedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			addresses: map[string]*uentity.Address{
				"user-id": {
					ID:        "address-id",
					UserID:    "user-id",
					IsDefault: true,
					AddressRevision: uentity.AddressRevision{
						ID:             1,
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
				},
			},
			expect: Users{
				{
					User: response.User{
						ID:            "user-id",
						Status:        int32(UserStatusVerified),
						Registered:    true,
						Username:      "username",
						AccountID:     "account-id",
						Lastname:      "&.",
						Firstname:     "購入者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "こうにゅうしゃ",
						Email:         "test-user@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
					address: Address{
						Address: response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							LastnameKana:   "あんどどっと",
							FirstnameKana:  "こうにゅうしゃ",
							PostalCode:     "1000014",
							Prefecture:     "東京都",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						revisionID: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewUsers(tt.users, tt.addresses))
		})
	}
}

func TestUsers_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  Users
		expect map[string]*User
	}{
		{
			name: "success",
			users: Users{
				{
					User: response.User{
						ID:            "user-id",
						Status:        int32(UserStatusGuest),
						Registered:    false,
						Username:      "",
						AccountID:     "",
						Lastname:      "&.",
						Firstname:     "購入者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "こうにゅうしゃ",
						Email:         "test-user@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
					address: Address{
						Address: response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							LastnameKana:   "あんどどっと",
							FirstnameKana:  "こうにゅうしゃ",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						revisionID: 1,
					},
				},
			},
			expect: map[string]*User{
				"user-id": {
					User: response.User{
						ID:            "user-id",
						Status:        int32(UserStatusGuest),
						Registered:    false,
						Username:      "",
						AccountID:     "",
						Lastname:      "&.",
						Firstname:     "購入者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "こうにゅうしゃ",
						Email:         "test-user@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
					address: Address{
						Address: response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							LastnameKana:   "あんどどっと",
							FirstnameKana:  "こうにゅうしゃ",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						revisionID: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.users.Map())
		})
	}
}

func TestUsers_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  Users
		expect []*response.User
	}{
		{
			name: "success",
			users: Users{
				{
					User: response.User{
						ID:            "user-id",
						Status:        int32(UserStatusGuest),
						Registered:    false,
						Username:      "",
						AccountID:     "",
						Lastname:      "&.",
						Firstname:     "購入者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "こうにゅうしゃ",
						Email:         "test-user@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
					address: Address{
						Address: response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							LastnameKana:   "あんどどっと",
							FirstnameKana:  "こうにゅうしゃ",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						revisionID: 1,
					},
				},
			},
			expect: []*response.User{
				{
					ID:            "user-id",
					Status:        int32(UserStatusGuest),
					Registered:    false,
					Username:      "",
					AccountID:     "",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.users.Response())
		})
	}
}

func TestUserToList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *User
		order  *sentity.AggregatedOrder
		expect *UserToList
	}{
		{
			name: "success",
			user: &User{
				User: response.User{
					ID:            "user-id",
					Status:        int32(UserStatusVerified),
					Registered:    true,
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				address: Address{
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					revisionID: 1,
				},
			},
			order: &sentity.AggregatedOrder{
				UserID:     "user-id",
				OrderCount: 2,
				Subtotal:   3000,
				Discount:   0,
			},
			expect: &UserToList{
				UserToList: response.UserToList{
					ID:                "user-id",
					Lastname:          "&.",
					Firstname:         "購入者",
					Email:             "test-user@and-period.jp",
					Status:            int32(UserStatusVerified),
					Registered:        true,
					PrefectureCode:    13,
					City:              "千代田区",
					PaymentTotalCount: 2,
				},
			},
		},
		{
			name: "success without order",
			user: &User{
				User: response.User{
					ID:            "user-id",
					Status:        int32(UserStatusVerified),
					Registered:    true,
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "購入者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "こうにゅうしゃ",
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				address: Address{
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
					revisionID: 1,
				},
			},
			order: nil,
			expect: &UserToList{
				UserToList: response.UserToList{
					ID:                "user-id",
					Lastname:          "&.",
					Firstname:         "購入者",
					Email:             "test-user@and-period.jp",
					Status:            int32(UserStatusVerified),
					Registered:        true,
					PrefectureCode:    13,
					City:              "千代田区",
					PaymentTotalCount: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewUserToList(tt.user, tt.order))
		})
	}
}

func TestUserToList_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *UserToList
		expect *response.UserToList
	}{
		{
			name: "success",
			user: &UserToList{
				UserToList: response.UserToList{
					ID:                "user-id",
					Lastname:          "&.",
					Firstname:         "購入者",
					Email:             "test-user@and-period.jp",
					Status:            int32(UserStatusGuest),
					Registered:        true,
					PrefectureCode:    13,
					City:              "千代田区",
					PaymentTotalCount: 2,
				},
			},
			expect: &response.UserToList{
				ID:                "user-id",
				Lastname:          "&.",
				Firstname:         "購入者",
				Email:             "test-user@and-period.jp",
				Status:            int32(UserStatusGuest),
				Registered:        true,
				PrefectureCode:    13,
				City:              "千代田区",
				PaymentTotalCount: 2,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Response())
		})
	}
}

func TestUsersToList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  Users
		orders map[string]*sentity.AggregatedOrder
		expect UsersToList
	}{
		{
			name: "success",
			users: Users{
				{
					User: response.User{
						ID:            "user-id",
						Status:        int32(UserStatusGuest),
						Registered:    false,
						Username:      "",
						AccountID:     "",
						Lastname:      "&.",
						Firstname:     "購入者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "こうにゅうしゃ",
						Email:         "test-user@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
					address: Address{
						Address: response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							LastnameKana:   "あんどどっと",
							FirstnameKana:  "こうにゅうしゃ",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						revisionID: 1,
					},
				},
			},
			orders: map[string]*sentity.AggregatedOrder{
				"user-id": {
					UserID:     "user-id",
					OrderCount: 2,
					Subtotal:   3000,
					Discount:   0,
				},
			},
			expect: UsersToList{
				{
					UserToList: response.UserToList{
						ID:                "user-id",
						Lastname:          "&.",
						Firstname:         "購入者",
						Email:             "test-user@and-period.jp",
						Status:            int32(UserStatusGuest),
						Registered:        false,
						PrefectureCode:    13,
						City:              "千代田区",
						PaymentTotalCount: 2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewUsersToList(tt.users, tt.orders))
		})
	}
}

func TestUsersToList_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  UsersToList
		expect []*response.UserToList
	}{
		{
			name: "success",
			users: UsersToList{
				{
					UserToList: response.UserToList{
						ID:                "user-id",
						Lastname:          "&.",
						Firstname:         "購入者",
						Email:             "test-user@and-period.jp",
						Status:            int32(UserStatusGuest),
						Registered:        true,
						PrefectureCode:    13,
						City:              "千代田区",
						PaymentTotalCount: 2,
					},
				},
			},
			expect: []*response.UserToList{
				{
					ID:                "user-id",
					Lastname:          "&.",
					Firstname:         "購入者",
					Email:             "test-user@and-period.jp",
					Status:            int32(UserStatusGuest),
					Registered:        true,
					PrefectureCode:    13,
					City:              "千代田区",
					PaymentTotalCount: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.users.Response())
		})
	}
}
