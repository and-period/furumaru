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
				Customer: uentity.Customer{
					UserID:         "user-id",
					Lastname:       "&.",
					Firstname:      "スタッフ",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "すたっふ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				Member: uentity.Member{
					UserID:       "user-id",
					AccountID:    "account-id",
					CognitoID:    "cognito-id",
					Username:     "username",
					ProviderType: uentity.ProviderTypeEmail,
					Email:        "test-user@and-period.jp",
					PhoneNumber:  "+819012345678",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
					CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
					VerifiedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
					ID:         "user-id",
					Registered: true,
					Email:      "test-user@and-period.jp",
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
				},
			},
		},
		{
			name: "success guest",
			user: &entity.User{
				ID:         "user-id",
				Registered: false,
				Customer: uentity.Customer{
					UserID:         "user-id",
					Lastname:       "&.",
					Firstname:      "スタッフ",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "すたっふ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
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
					ID:         "user-id",
					Registered: false,
					Email:      "test-user@and-period.jp",
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
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
					ID:         "user-id",
					Registered: true,
					Email:      "test-user@and-period.jp",
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
				},
			},
			expect: &response.User{
				ID:         "user-id",
				Registered: true,
				Email:      "test-user@and-period.jp",
				CreatedAt:  1640962800,
				UpdatedAt:  1640962800,
				Address: &response.Address{
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
				},
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
					Customer: uentity.Customer{
						UserID:         "user-id",
						Lastname:       "&.",
						Firstname:      "スタッフ",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "すたっふ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					Member: uentity.Member{
						UserID:       "user-id",
						AccountID:    "account-id",
						CognitoID:    "cognito-id",
						Username:     "username",
						ProviderType: uentity.ProviderTypeEmail,
						Email:        "test-user@and-period.jp",
						PhoneNumber:  "+819012345678",
						ThumbnailURL: "https://and-period.jp/thumbnail.png",
						CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
						VerifiedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
						ID:         "user-id",
						Registered: true,
						Email:      "test-user@and-period.jp",
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
						Address: &response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
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
						ID:         "user-id",
						Registered: false,
						Email:      "test-user@and-period.jp",
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
						Address: &response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
					},
				},
			},
			expect: map[string]*User{
				"user-id": {
					User: response.User{
						ID:         "user-id",
						Registered: false,
						Email:      "test-user@and-period.jp",
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
						Address: &response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
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
						ID:         "user-id",
						Registered: false,
						Email:      "test-user@and-period.jp",
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
						Address: &response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
					},
				},
			},
			expect: []*response.User{
				{
					ID:         "user-id",
					Registered: false,
					Email:      "test-user@and-period.jp",
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
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

func TestUserSummary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *User
		order  *sentity.AggregatedOrder
		expect *UserSummary
	}{
		{
			name: "success",
			user: &User{
				User: response.User{
					ID:         "user-id",
					Registered: true,
					Email:      "test-user@and-period.jp",
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
				},
			},
			order: &sentity.AggregatedOrder{
				UserID:     "user-id",
				OrderCount: 2,
				Subtotal:   3000,
				Discount:   0,
			},
			expect: &UserSummary{
				UserSummary: response.UserSummary{
					ID:             "user-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					Registered:     true,
					PrefectureCode: 13,
					City:           "千代田区",
					TotalOrder:     2,
					TotalAmount:    3000,
				},
			},
		},
		{
			name: "success without order",
			user: &User{
				User: response.User{
					ID:         "user-id",
					Registered: false,
					Email:      "test-user@and-period.jp",
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "+819012345678",
					},
				},
			},
			order: nil,
			expect: &UserSummary{
				UserSummary: response.UserSummary{
					ID:             "user-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					Registered:     false,
					PrefectureCode: 13,
					City:           "千代田区",
					TotalOrder:     0,
					TotalAmount:    0,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewUserSummary(tt.user, tt.order))
		})
	}
}

func TestUserSummary_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *UserSummary
		expect *response.UserSummary
	}{
		{
			name: "success",
			user: &UserSummary{
				UserSummary: response.UserSummary{
					ID:             "user-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					Registered:     true,
					PrefectureCode: 13,
					City:           "千代田区",
					TotalOrder:     2,
					TotalAmount:    3000,
				},
			},
			expect: &response.UserSummary{
				ID:             "user-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				Registered:     true,
				PrefectureCode: 13,
				City:           "千代田区",
				TotalOrder:     2,
				TotalAmount:    3000,
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

func TestUserSummaries(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  Users
		orders map[string]*sentity.AggregatedOrder
		expect UserSummaries
	}{
		{
			name: "success",
			users: Users{
				{
					User: response.User{
						ID:         "user-id",
						Registered: false,
						Email:      "test-user@and-period.jp",
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
						Address: &response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
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
			expect: UserSummaries{
				{
					UserSummary: response.UserSummary{
						ID:             "user-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						Registered:     false,
						PrefectureCode: 13,
						City:           "千代田区",
						TotalOrder:     2,
						TotalAmount:    3000,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewUserSummaries(tt.users, tt.orders))
		})
	}
}

func TestUserSummaries_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  UserSummaries
		expect []*response.UserSummary
	}{
		{
			name: "success",
			users: UserSummaries{
				{
					UserSummary: response.UserSummary{
						ID:             "user-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						Registered:     true,
						PrefectureCode: 13,
						City:           "千代田区",
						TotalOrder:     2,
						TotalAmount:    3000,
					},
				},
			},
			expect: []*response.UserSummary{
				{
					ID:             "user-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					Registered:     true,
					PrefectureCode: 13,
					City:           "千代田区",
					TotalOrder:     2,
					TotalAmount:    3000,
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
