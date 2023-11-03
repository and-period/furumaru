package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		address *entity.Address
		expect  *Address
	}{
		{
			name: "success",
			address: &entity.Address{
				ID:             "address-id",
				UserID:         "user-id",
				Hash:           "789ef22a79a364f95c66a3d3b1fda213c1316a6c7f8b6306b493d8c46d2dce75",
				IsDefault:      true,
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
			expect: &Address{
				Address: response.Address{
					ID:           "address-id",
					IsDefault:    true,
					Lastname:     "&.",
					Firstname:    "購入者",
					PostalCode:   "1000014",
					Prefecture:   "東京都",
					City:         "千代田区",
					AddressLine1: "永田町1-7-1",
					AddressLine2: "",
					PhoneNumber:  "+819012345678",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAddress(tt.address)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAddress_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		address *Address
		expect  *response.Address
	}{
		{
			name: "success",
			address: &Address{
				Address: response.Address{
					ID:           "address-id",
					IsDefault:    true,
					Lastname:     "&.",
					Firstname:    "購入者",
					PostalCode:   "1000014",
					Prefecture:   "東京都",
					City:         "千代田区",
					AddressLine1: "永田町1-7-1",
					AddressLine2: "",
					PhoneNumber:  "+819012345678",
				},
			},
			expect: &response.Address{
				ID:           "address-id",
				IsDefault:    true,
				Lastname:     "&.",
				Firstname:    "購入者",
				PostalCode:   "1000014",
				Prefecture:   "東京都",
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
				PhoneNumber:  "+819012345678",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.address.Response())
		})
	}
}

func TestAddresses(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses entity.Addresses
		expect    Addresses
	}{
		{
			name: "success",
			addresses: entity.Addresses{
				{
					ID:             "address-id",
					UserID:         "user-id",
					Hash:           "789ef22a79a364f95c66a3d3b1fda213c1316a6c7f8b6306b493d8c46d2dce75",
					IsDefault:      true,
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
			expect: Addresses{
				{
					Address: response.Address{
						ID:           "address-id",
						IsDefault:    true,
						Lastname:     "&.",
						Firstname:    "購入者",
						PostalCode:   "1000014",
						Prefecture:   "東京都",
						City:         "千代田区",
						AddressLine1: "永田町1-7-1",
						AddressLine2: "",
						PhoneNumber:  "+819012345678",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAddresses(tt.addresses)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAddresses_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expect    []*response.Address
	}{
		{
			name: "success",
			addresses: Addresses{
				{
					Address: response.Address{
						ID:           "address-id",
						IsDefault:    true,
						Lastname:     "&.",
						Firstname:    "購入者",
						PostalCode:   "1000014",
						Prefecture:   "東京都",
						City:         "千代田区",
						AddressLine1: "永田町1-7-1",
						AddressLine2: "",
						PhoneNumber:  "+819012345678",
					},
				},
			},
			expect: []*response.Address{
				{
					ID:           "address-id",
					IsDefault:    true,
					Lastname:     "&.",
					Firstname:    "購入者",
					PostalCode:   "1000014",
					Prefecture:   "東京都",
					City:         "千代田区",
					AddressLine1: "永田町1-7-1",
					AddressLine2: "",
					PhoneNumber:  "+819012345678",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.addresses.Response())
		})
	}
}
