package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewUserParams
		expect *User
	}{
		{
			name: "success with member",
			params: &NewUserParams{
				Registered:   true,
				CognitoID:    "cognito-id",
				ProviderType: ProviderTypeEmail,
				Email:        "test-user@and-period.jp",
				PhoneNumber:  "+810000000000",
			},
			expect: &User{
				Registered: true,
				Member: Member{
					CognitoID:    "cognito-id",
					ProviderType: ProviderTypeEmail,
					Email:        "test-user@and-period.jp",
					PhoneNumber:  "+810000000000",
				},
			},
		},
		{
			name: "success with guest",
			params: &NewUserParams{
				Registered:   false,
				CognitoID:    "cognito-id",
				ProviderType: ProviderTypeEmail,
				Email:        "test-user@and-period.jp",
				PhoneNumber:  "+810000000000",
			},
			expect: &User{
				Registered: false,
				Guest: Guest{
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+810000000000",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewUser(tt.params)
			actual.ID = ""              // ignore
			actual.Member.UserID = ""   // ignore
			actual.Guest.UserID = ""    // ignore
			actual.Customer.UserID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestUser_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		user   *User
		expect string
	}{
		{
			name: "success",
			user: &User{
				Customer: Customer{
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
			},
			expect: "&. スタッフ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Name())
		})
	}
}

func TestUser_Email(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *User
		expect string
	}{
		{
			name: "success member",
			user: &User{
				ID:         "user-id",
				Registered: true,
				Customer: Customer{
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
				Member: Member{
					UserID:       "user-id",
					AccountID:    "account-id",
					CognitoID:    "cognito-id",
					Username:     "username",
					ProviderType: ProviderTypeEmail,
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
			expect: "test-user@and-period.jp",
		},
		{
			name: "success guest",
			user: &User{
				ID:         "user-id",
				Registered: false,
				Customer: Customer{
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
				Guest: Guest{
					UserID:      "user-id",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: "test-user@and-period.jp",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Email())
		})
	}
}

func TestUser_PhoneNumber(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *User
		expect string
	}{
		{
			name: "success member",
			user: &User{
				ID:         "user-id",
				Registered: true,
				Customer: Customer{
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
				Member: Member{
					UserID:       "user-id",
					AccountID:    "account-id",
					CognitoID:    "cognito-id",
					Username:     "username",
					ProviderType: ProviderTypeEmail,
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
			expect: "+819012345678",
		},
		{
			name: "success guest",
			user: &User{
				ID:         "user-id",
				Registered: false,
				Customer: Customer{
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
				Guest: Guest{
					UserID:      "user-id",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: "+819012345678",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.PhoneNumber())
		})
	}
}

func TestUser_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		user     *User
		customer *Customer
		member   *Member
		guest    *Guest
		expect   *User
	}{
		{
			name: "success",
			user: &User{},
			customer: &Customer{
				Lastname:  "&.",
				Firstname: "スタッフ",
			},
			member: &Member{
				UserID: "user-id",
			},
			guest: &Guest{
				UserID: "user-id",
			},
			expect: &User{
				Customer: Customer{
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
				Member: Member{
					UserID: "user-id",
				},
				Guest: Guest{
					UserID: "user-id",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.user.Fill(tt.customer, tt.member, tt.guest)
			assert.Equal(t, tt.expect, tt.user)
		})
	}
}

func TestUsers_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  Users
		expect []string
	}{
		{
			name: "success",
			users: Users{
				{ID: "user-id01"},
				{ID: "user-id02"},
			},
			expect: []string{
				"user-id01",
				"user-id02",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.users.IDs())
		})
	}
}

func TestUsers_GroupByRegistered(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  Users
		expect map[bool]Users
	}{
		{
			name: "success",
			users: Users{
				{
					ID:         "user-id01",
					Registered: true,
					Member: Member{
						UserID: "user-id01",
					},
				},
				{
					ID:         "user-id02",
					Registered: false,
					Guest: Guest{
						UserID: "user-id02",
					},
				},
			},
			expect: map[bool]Users{
				true: {
					{
						ID:         "user-id01",
						Registered: true,
						Member: Member{
							UserID: "user-id01",
						},
					},
				},
				false: {
					{
						ID:         "user-id02",
						Registered: false,
						Guest: Guest{
							UserID: "user-id02",
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
			assert.Equal(t, tt.expect, tt.users.GroupByRegistered())
		})
	}
}
