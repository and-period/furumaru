package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		user   *entity.User
		expect *AuthUser
	}{
		{
			name: "success",
			user: &entity.User{
				ID: "user-id",
				FacilityUser: entity.FacilityUser{
					Firstname:     "太郎",
					Lastname:      "山田",
					FirstnameKana: "たろう",
					LastnameKana:  "やまだ",
					Email:         "test@example.com",
					PhoneNumber:   "090-1234-5678",
					LastCheckInAt: time.Unix(1640995200, 0), // 2022-01-01 00:00:00 UTC
				},
			},
			expect: &AuthUser{
				AuthUser: types.AuthUser{
					ID:            "user-id",
					Firstname:     "太郎",
					Lastname:      "山田",
					FirstnameKana: "たろう",
					LastnameKana:  "やまだ",
					Email:         "test@example.com",
					PhoneNumber:   "090-1234-5678",
					LastCheckInAt: 1640995200,
				},
			},
		},
		{
			name: "success with empty facility user fields",
			user: &entity.User{
				ID: "user-id",
				FacilityUser: entity.FacilityUser{
					Firstname:     "",
					Lastname:      "",
					FirstnameKana: "",
					LastnameKana:  "",
					Email:         "",
					PhoneNumber:   "",
					LastCheckInAt: time.Unix(0, 0),
				},
			},
			expect: &AuthUser{
				AuthUser: types.AuthUser{
					ID:            "user-id",
					Firstname:     "",
					Lastname:      "",
					FirstnameKana: "",
					LastnameKana:  "",
					Email:         "",
					PhoneNumber:   "",
					LastCheckInAt: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAuthUser(tt.user)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAuthUser_Response(t *testing.T) {
	t.Parallel()

	authUser := &AuthUser{
		AuthUser: types.AuthUser{
			ID:            "user-id",
			Firstname:     "太郎",
			Lastname:      "山田",
			FirstnameKana: "たろう",
			LastnameKana:  "やまだ",
			Email:         "test@example.com",
			PhoneNumber:   "090-1234-5678",
			LastCheckInAt: 1640995200,
		},
	}

	expected := &types.AuthUser{
		ID:            "user-id",
		Firstname:     "太郎",
		Lastname:      "山田",
		FirstnameKana: "たろう",
		LastnameKana:  "やまだ",
		Email:         "test@example.com",
		PhoneNumber:   "090-1234-5678",
		LastCheckInAt: 1640995200,
	}

	assert.Equal(t, expected, authUser.Response())
}
