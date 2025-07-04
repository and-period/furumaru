package entity

import (
	"testing"
	"time"

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
				Registered:    true,
				CognitoID:     "cognito-id",
				Username:      "username",
				AccountID:     "account-id",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				ProviderType:  UserAuthProviderTypeEmail,
				Email:         "test-user@and-period.jp",
				PhoneNumber:   "+810000000000",
			},
			expect: &User{
				Registered: true,
				Member: Member{
					CognitoID:     "cognito-id",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					ProviderType:  UserAuthProviderTypeEmail,
					Email:         "test-user@and-period.jp",
					PhoneNumber:   "+810000000000",
				},
			},
		},
		{
			name: "success with guest",
			params: &NewUserParams{
				Registered:    false,
				CognitoID:     "cognito-id",
				Username:      "username",
				AccountID:     "account-id",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				ProviderType:  UserAuthProviderTypeEmail,
				Email:         "test-user@and-period.jp",
				PhoneNumber:   "+810000000000",
			},
			expect: &User{
				Registered: false,
				Guest: Guest{
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					Email:         "test-user@and-period.jp",
				},
			},
		},
	}

	for _, tt := range tests {
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
			name: "success member",
			user: &User{
				ID:         "user-id",
				Registered: true,
				Member: Member{
					UserID:        "user-id",
					AccountID:     "account-id",
					CognitoID:     "cognito-id",
					Username:      "username",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					ProviderType:  UserAuthProviderTypeEmail,
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
			expect: "&. 利用者",
		},
		{
			name: "success guest",
			user: &User{
				ID:         "user-id",
				Registered: false,
				Guest: Guest{
					UserID:    "user-id",
					Email:     "test-user@and-period.jp",
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: "ゲスト",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Name())
		})
	}
}

func TestUser_Username(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *User
		expect string
	}{
		{
			name: "success member",
			user: &User{
				Registered: true,
				Member: Member{
					Username: "username",
				},
			},
			expect: "username",
		},
		{
			name: "success guest",
			user: &User{
				Registered: false,
			},
			expect: "名無しさん",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Username())
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
				Member: Member{
					UserID:       "user-id",
					AccountID:    "account-id",
					CognitoID:    "cognito-id",
					Username:     "username",
					ProviderType: UserAuthProviderTypeEmail,
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
				Guest: Guest{
					UserID:    "user-id",
					Email:     "test-user@and-period.jp",
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: "test-user@and-period.jp",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Email())
		})
	}
}

func TestUser_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		user   *User
		member *Member
		guest  *Guest
		expect *User
	}{
		{
			name: "success",
			user: &User{},
			member: &Member{
				UserID: "user-id",
			},
			guest: &Guest{
				UserID: "user-id",
			},
			expect: &User{
				Member: Member{
					UserID: "user-id",
				},
				Guest: Guest{
					UserID: "user-id",
				},
				Status: UserStatusGuest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.user.Fill(tt.member, tt.guest)
			assert.Equal(t, tt.expect, tt.user)
		})
	}
}

func TestUser_SetStatus(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		user   *User
		expect *User
	}{
		{
			name: "guest",
			user: &User{
				Registered: false,
			},
			expect: &User{
				Registered: false,
				Status:     UserStatusGuest,
			},
		},
		{
			name: "provisional",
			user: &User{
				Member: Member{
					VerifiedAt: time.Time{},
				},
				Registered: true,
			},
			expect: &User{
				Registered: true,
				Status:     UserStatusProvisional,
			},
		},
		{
			name: "verified",
			user: &User{
				Member: Member{
					VerifiedAt: now,
				},
				Registered: true,
			},
			expect: &User{
				Member: Member{
					VerifiedAt: now,
				},
				Registered: true,
				Status:     UserStatusVerified,
			},
		},
		{
			name:   "empty",
			user:   nil,
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.user.SetStatus()
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.users.IDs())
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
			expect: map[string]*User{
				"user-id01": {
					ID:         "user-id01",
					Registered: true,
					Member: Member{
						UserID: "user-id01",
					},
				},
				"user-id02": {
					ID:         "user-id02",
					Registered: false,
					Guest: Guest{
						UserID: "user-id02",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.users.Map())
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.users.GroupByRegistered())
		})
	}
}

func TestUsers_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		users   Users
		members map[string]*Member
		guests  map[string]*Guest
		expect  Users
	}{
		{
			name: "success",
			users: Users{
				{ID: "user-id01"},
				{ID: "user-id02"},
			},
			members: map[string]*Member{
				"user-id01": {
					UserID: "user-id01",
				},
			},
			guests: map[string]*Guest{
				"user-id01": {
					UserID: "user-id01",
				},
			},
			expect: Users{
				{
					ID:     "user-id01",
					Member: Member{UserID: "user-id01"},
					Guest:  Guest{UserID: "user-id01"},
					Status: UserStatusGuest,
				},
				{
					ID:     "user-id02",
					Member: Member{},
					Guest:  Guest{},
					Status: UserStatusGuest,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.users.Fill(tt.members, tt.guests)
			assert.Equal(t, tt.expect, tt.users)
		})
	}
}
