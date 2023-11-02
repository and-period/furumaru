package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewAddressParams
		expect *Address
		hasErr bool
	}{
		{
			name: "success",
			params: &NewAddressParams{
				UserID:       "user-id",
				IsDefault:    true,
				Lastname:     "&.",
				Firstname:    "購入者",
				PostalCode:   "1000014",
				Prefecture:   13,
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
				PhoneNumber:  "+819012345678",
			},
			expect: &Address{
				UserID:       "user-id",
				Hash:         "c1f66591133a1a70cc6b29f21ede4389efe6864bb7ade2e17f734471352df1a9",
				IsDefault:    true,
				Lastname:     "&.",
				Firstname:    "購入者",
				PostalCode:   "1000014",
				Prefecture:   13,
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
				PhoneNumber:  "+819012345678",
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture",
			params: &NewAddressParams{
				UserID:       "user-id",
				IsDefault:    true,
				Lastname:     "&.",
				Firstname:    "購入者",
				PostalCode:   "1000014",
				Prefecture:   0,
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
				PhoneNumber:  "+819012345678",
			},
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewAddress(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			if actual != nil {
				actual.ID = "" // ignore
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAddressHash(t *testing.T) {
	t.Parallel()
	type args struct {
		userID       string
		postalCode   string
		addressLine1 string
		addressLine2 string
	}
	tests := []struct {
		name   string
		params *NewAddressParams
		expect string
	}{
		{
			name: "success",
			params: &NewAddressParams{
				UserID:       "user-id",
				IsDefault:    true,
				Lastname:     "&.",
				Firstname:    "購入者",
				PostalCode:   "1000014",
				Prefecture:   0,
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
				PhoneNumber:  "+819012345678",
			},
			expect: "c1f66591133a1a70cc6b29f21ede4389efe6864bb7ade2e17f734471352df1a9",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAddressHash(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
