package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacilityUser_Name(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		facilityUser *FacilityUser
		expect       string
	}{
		{
			name: "success",
			facilityUser: &FacilityUser{
				UserID:        "user-id",
				ExternalID:    "external-id",
				ProducerID:    "producer-id",
				Lastname:      "&.",
				Firstname:     "施設利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "しせつりようしゃ",
				ProviderType:  UserAuthProviderTypeEmail,
				Email:         "test@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expect: "&. 施設利用者",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.facilityUser.Name())
		})
	}
}

func TestFacilityUsers_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		users  FacilityUsers
		expect map[string]*FacilityUser
	}{
		{
			name: "success",
			users: FacilityUsers{
				{
					UserID: "user-id01",
					Email:  "test-user01@and-period.jp",
				},
				{
					UserID: "user-id02",
					Email:  "test-user02@and-period.jp",
				},
			},
			expect: map[string]*FacilityUser{
				"user-id01": {
					UserID: "user-id01",
					Email:  "test-user01@and-period.jp",
				},
				"user-id02": {
					UserID: "user-id02",
					Email:  "test-user02@and-period.jp",
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
