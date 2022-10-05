package entity

import (
	"testing"

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
			actual.ID = ""            // ignore
			actual.Member.UserID = "" // ignore
			actual.Guest.UserID = ""  // ignore
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
